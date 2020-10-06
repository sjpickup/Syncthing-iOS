// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package discover

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/syncthing/protocol"
	"../../internal/beacon"
	"../../internal/events"
	"../../internal/sync"
)

type Discoverer struct {
	myID            protocol.DeviceID
	listenAddrs     []string
	localBcastIntv  time.Duration
	localBcastStart time.Time
	cacheLifetime   time.Duration
	negCacheCutoff  time.Duration
	beacons         []beacon.Interface
	extPort         uint16
	localBcastTick  <-chan time.Time
	forcedBcastTick chan time.Time

	registryLock sync.RWMutex
	registry     map[protocol.DeviceID][]CacheEntry
	lastLookup   map[protocol.DeviceID]time.Time

	clients []Client
	mut     sync.RWMutex
}

type CacheEntry struct {
	Address string
	Seen    time.Time
}

var (
	ErrIncorrectMagic = errors.New("incorrect magic number")
)

func NewDiscoverer(id protocol.DeviceID, addresses []string) *Discoverer {
	return &Discoverer{
		myID:           id,
		listenAddrs:    addresses,
		localBcastIntv: 30 * time.Second,
		cacheLifetime:  5 * time.Minute,
		negCacheCutoff: 3 * time.Minute,
		registry:       make(map[protocol.DeviceID][]CacheEntry),
		lastLookup:     make(map[protocol.DeviceID]time.Time),
		registryLock:   sync.NewRWMutex(),
		mut:            sync.NewRWMutex(),
	}
}

func (d *Discoverer) StartLocal(localPort int, localMCAddr string) {
	if localPort > 0 {
		d.startLocalIPv4Broadcasts(localPort)
	}

	if len(localMCAddr) > 0 {
		d.startLocalIPv6Multicasts(localMCAddr)
	}

	if len(d.beacons) == 0 {
		l.Warnln("Local discovery unavailable")
		return
	}

	d.localBcastTick = time.Tick(d.localBcastIntv)
	d.forcedBcastTick = make(chan time.Time)
	d.localBcastStart = time.Now()
	go d.sendLocalAnnouncements()
}

func (d *Discoverer) startLocalIPv4Broadcasts(localPort int) {
	bb := beacon.NewBroadcast(localPort)
	d.beacons = append(d.beacons, bb)
	go d.recvAnnouncements(bb)
	bb.ServeBackground()
}

func (d *Discoverer) startLocalIPv6Multicasts(localMCAddr string) {
	intfs, err := net.Interfaces()
	if err != nil {
		if debug {
			l.Debugln("discover: interfaces:", err)
		}
		l.Infoln("Local discovery over IPv6 unavailable")
		return
	}

	v6Intfs := 0
	for _, intf := range intfs {
		// Interface flags seem to always be 0 on Windows
		if runtime.GOOS != "windows" && (intf.Flags&net.FlagUp == 0 || intf.Flags&net.FlagMulticast == 0) {
			continue
		}

		mb, err := beacon.NewMulticast(localMCAddr, intf.Name)
		if err != nil {
			if debug {
				l.Debugln("discover: Start local v6:", err)
			}
			continue
		}

		d.beacons = append(d.beacons, mb)
		go d.recvAnnouncements(mb)
		v6Intfs++
	}

	if v6Intfs == 0 {
		l.Infoln("Local discovery over IPv6 unavailable")
	}
}

func (d *Discoverer) StartGlobal(servers []string, extPort uint16) {
	d.mut.Lock()
	defer d.mut.Unlock()

	if len(d.clients) > 0 {
		d.stopGlobal()
	}

	d.extPort = extPort
	pkt := d.announcementPkt()
	wg := sync.NewWaitGroup()
	clients := make(chan Client, len(servers))
	for _, address := range servers {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			client, err := New(addr, pkt)
			if err != nil {
				l.Infoln("Error creating discovery client", addr, err)
				return
			}
			clients <- client
		}(address)
	}

	wg.Wait()
	close(clients)

	for client := range clients {
		d.clients = append(d.clients, client)
	}
}

func (d *Discoverer) StopGlobal() {
	d.mut.Lock()
	defer d.mut.Unlock()
	d.stopGlobal()
}

func (d *Discoverer) stopGlobal() {
	for _, client := range d.clients {
		client.Stop()
	}
	d.clients = []Client{}
}

func (d *Discoverer) ExtAnnounceOK() map[string]bool {
	d.mut.RLock()
	defer d.mut.RUnlock()

	ret := make(map[string]bool)
	for _, client := range d.clients {
		ret[client.Address()] = client.StatusOK()
	}
	return ret
}

func (d *Discoverer) Lookup(device protocol.DeviceID) []string {
	d.registryLock.RLock()
	cached := d.filterCached(d.registry[device])
	lastLookup := d.lastLookup[device]
	d.registryLock.RUnlock()

	d.mut.RLock()
	defer d.mut.RUnlock()

	if len(cached) > 0 {
		// There are cached address entries.
		addrs := make([]string, len(cached))
		for i := range cached {
			addrs[i] = cached[i].Address
		}
		return addrs
	}

	if time.Since(lastLookup) < d.negCacheCutoff {
		// We have recently tried to lookup this address and failed. Lets
		// chill for a while.
		return nil
	}

	if len(d.clients) != 0 && time.Since(d.localBcastStart) > d.localBcastIntv {
		// Only perform external lookups if we have at least one external
		// server client and one local announcement interval has passed. This is
		// to avoid finding local peers on their remote address at startup.
		results := make(chan []string, len(d.clients))
		wg := sync.NewWaitGroup()
		for _, client := range d.clients {
			wg.Add(1)
			go func(c Client) {
				defer wg.Done()
				results <- c.Lookup(device)
			}(client)
		}

		wg.Wait()
		close(results)

		cached := []CacheEntry{}
		seen := make(map[string]struct{})
		now := time.Now()

		var addrs []string
		for result := range results {
			for _, addr := range result {
				_, ok := seen[addr]
				if !ok {
					cached = append(cached, CacheEntry{
						Address: addr,
						Seen:    now,
					})
					seen[addr] = struct{}{}
					addrs = append(addrs, addr)
				}
			}
		}

		d.registryLock.Lock()
		d.registry[device] = cached
		d.lastLookup[device] = time.Now()
		d.registryLock.Unlock()

		return addrs
	}

	return nil
}

func (d *Discoverer) Hint(device string, addrs []string) {
	resAddrs := resolveAddrs(addrs)
	var id protocol.DeviceID
	id.UnmarshalText([]byte(device))
	d.registerDevice(nil, Device{
		Addresses: resAddrs,
		ID:        id[:],
	})
}

func (d *Discoverer) All() map[protocol.DeviceID][]CacheEntry {
	d.registryLock.RLock()
	devices := make(map[protocol.DeviceID][]CacheEntry, len(d.registry))
	for device, addrs := range d.registry {
		addrsCopy := make([]CacheEntry, len(addrs))
		copy(addrsCopy, addrs)
		devices[device] = addrsCopy
	}
	d.registryLock.RUnlock()
	return devices
}

func (d *Discoverer) announcementPkt() *Announce {
	var addrs []Address
	if d.extPort != 0 {
		addrs = []Address{{Port: d.extPort}}
	} else {
		for _, astr := range d.listenAddrs {
			addr, err := net.ResolveTCPAddr("tcp", astr)
			if err != nil {
				l.Warnln("discover: %v: not announcing %s", err, astr)
				continue
			} else if debug {
				l.Debugf("discover: resolved %s as %#v", astr, addr)
			}
			if len(addr.IP) == 0 || addr.IP.IsUnspecified() {
				addrs = append(addrs, Address{Port: uint16(addr.Port)})
			} else if bs := addr.IP.To4(); bs != nil {
				addrs = append(addrs, Address{IP: bs, Port: uint16(addr.Port)})
			} else if bs := addr.IP.To16(); bs != nil {
				addrs = append(addrs, Address{IP: bs, Port: uint16(addr.Port)})
			}
		}
	}
	return &Announce{
		Magic: AnnouncementMagic,
		This:  Device{d.myID[:], addrs},
	}
}

func (d *Discoverer) sendLocalAnnouncements() {
	var addrs = resolveAddrs(d.listenAddrs)

	var pkt = Announce{
		Magic: AnnouncementMagic,
		This:  Device{d.myID[:], addrs},
	}
	msg := pkt.MustMarshalXDR()

	for {
		for _, b := range d.beacons {
			b.Send(msg)
		}

		select {
		case <-d.localBcastTick:
		case <-d.forcedBcastTick:
		}
	}
}

func (d *Discoverer) recvAnnouncements(b beacon.Interface) {
	for {
		buf, addr := b.Recv()

		var pkt Announce
		err := pkt.UnmarshalXDR(buf)
		if err != nil && err != io.EOF {
			if debug {
				l.Debugf("discover: Failed to unmarshal local announcement from %s:\n%s", addr, hex.Dump(buf))
			}
			continue
		}

		if debug {
			l.Debugf("discover: Received local announcement from %s for %s", addr, protocol.DeviceIDFromBytes(pkt.This.ID))
		}

		var newDevice bool
		if bytes.Compare(pkt.This.ID, d.myID[:]) != 0 {
			newDevice = d.registerDevice(addr, pkt.This)
		}

		if newDevice {
			select {
			case d.forcedBcastTick <- time.Now():
			}
		}
	}
}

func (d *Discoverer) registerDevice(addr net.Addr, device Device) bool {
	var id protocol.DeviceID
	copy(id[:], device.ID)

	d.registryLock.Lock()
	defer d.registryLock.Unlock()

	current := d.filterCached(d.registry[id])

	orig := current

	for _, a := range device.Addresses {
		var deviceAddr string
		if len(a.IP) > 0 {
			deviceAddr = net.JoinHostPort(net.IP(a.IP).String(), strconv.Itoa(int(a.Port)))
		} else if addr != nil {
			ua := addr.(*net.UDPAddr)
			ua.Port = int(a.Port)
			deviceAddr = ua.String()
		}
		for i := range current {
			if current[i].Address == deviceAddr {
				current[i].Seen = time.Now()
				goto done
			}
		}
		current = append(current, CacheEntry{
			Address: deviceAddr,
			Seen:    time.Now(),
		})
	done:
	}

	if debug {
		l.Debugf("discover: Caching %s addresses: %v", id, current)
	}

	d.registry[id] = current

	if len(current) > len(orig) {
		addrs := make([]string, len(current))
		for i := range current {
			addrs[i] = current[i].Address
		}
		events.Default.Log(events.DeviceDiscovered, map[string]interface{}{
			"device": id.String(),
			"addrs":  addrs,
		})
	}

	return len(current) > len(orig)
}

func (d *Discoverer) filterCached(c []CacheEntry) []CacheEntry {
	for i := 0; i < len(c); {
		if ago := time.Since(c[i].Seen); ago > d.cacheLifetime {
			if debug {
				l.Debugf("discover: Removing cached address %s - seen %v ago", c[i].Address, ago)
			}
			c[i] = c[len(c)-1]
			c = c[:len(c)-1]
		} else {
			i++
		}
	}
	return c
}

func addrToAddr(addr *net.TCPAddr) Address {
	if len(addr.IP) == 0 || addr.IP.IsUnspecified() {
		return Address{Port: uint16(addr.Port)}
	} else if bs := addr.IP.To4(); bs != nil {
		return Address{IP: bs, Port: uint16(addr.Port)}
	} else if bs := addr.IP.To16(); bs != nil {
		return Address{IP: bs, Port: uint16(addr.Port)}
	}
	return Address{}
}

func resolveAddrs(addrs []string) []Address {
	var raddrs []Address
	for _, addrStr := range addrs {
		addrRes, err := net.ResolveTCPAddr("tcp", addrStr)
		if err != nil {
			continue
		}
		addr := addrToAddr(addrRes)
		if len(addr.IP) > 0 {
			raddrs = append(raddrs, addr)
		} else {
			raddrs = append(raddrs, Address{Port: addr.Port})
		}
	}
	return raddrs
}
