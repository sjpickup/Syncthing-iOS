// Created by cgo - DO NOT EDIT

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:9
package main
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:12

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:11
import (
												_ "runtime/cgo"
												"crypto/tls"
												"flag"
												"fmt"
												"io/ioutil"
												"log"
												"net"
												"net/http"
												_ "net/http/pprof"
												"net/url"
												"os"
												"path/filepath"
												"regexp"
												"runtime"
												"runtime/pprof"
												"strconv"
												"strings"
												"time"
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:33

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:32
	"github.com/calmh/logger"
												"github.com/juju/ratelimit"
												"github.com/syncthing/protocol"
												"../../internal/config"
												"../../internal/db"
												"../../internal/discover"
												"../../internal/events"
												"../../internal/model"
												"../../internal/osutil"
												"../../internal/symlinks"
												"../../internal/upgrade"
												"github.com/syndtr/goleveldb/leveldb"
												"github.com/syndtr/goleveldb/leveldb/errors"
												"github.com/syndtr/goleveldb/leveldb/opt"
												"github.com/thejerf/suture"
												"golang.org/x/crypto/bcrypt"
)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:57

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:56
var (
	Version		= "unknown-dev"
	BuildEnv	= "default"
	BuildStamp	= "0"
	BuildDate	time.Time
	BuildHost	= "unknown"
	BuildUser	= "unknown"
	IsRelease	bool
	IsBeta		bool
	LongVersion	string
)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:69

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:68
const (
	exitSuccess		= 0
	exitError		= 1
	exitNoUpgradeAvailable	= 2
	exitRestarting		= 3
	exitUpgrading		= 4
)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:77

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:76
const (
	bepProtocolName		= "bep/1.0"
	pingEventInterval	= time.Minute
)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:82

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:81
var l = logger.DefaultLogger
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:85

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:84
func init() {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:87

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:86
	if Version != "unknown-dev" {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:89

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:88
		exp := regexp.MustCompile(`^v\d+\.\d+\.\d+(-[a-z0-9]+)*(\+\d+-g[0-9a-f]+)?(-dirty)?$`)
													if !exp.MatchString(Version) {
			l.Fatalf("Invalid version string %q;\n\tdoes not match regexp %v", Version, exp)
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:101

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:100
	exp := regexp.MustCompile(`^v\d+\.\d+\.\d+(-[a-z]+[\d\.]+)?$`)
												IsRelease = exp.MatchString(Version)
												IsBeta = strings.Contains(Version, "-")
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:105

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:104
	stamp, _ := strconv.Atoi(BuildStamp)
												BuildDate = time.Unix(int64(stamp), 0)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:108

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:107
	date := BuildDate.UTC().Format("2006-01-02 15:04:05 MST")
												LongVersion = fmt.Sprintf("syncthing %s (%s %s-%s %s) %s@%s %s", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, BuildEnv, BuildUser, BuildHost, date)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:111

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:110
	if os.Getenv("STTRACE") != "" {
		logFlags = log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:116

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:115
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:118

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:117
var (
	cfg		*config.Wrapper
	myID		protocol.DeviceID
	confDir		string
	logFlags	= log.Ltime
	writeRateLimit	*ratelimit.Bucket
	readRateLimit	*ratelimit.Bucket
	stop		= make(chan int)
	discoverer	*discover.Discoverer
	cert		tls.Certificate
	lans		[]*net.IPNet
)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:131

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:130
const (
	usage		= "syncthing [options]"
	extraUsage	= `
The default configuration directory is:

  %s


The -logflags value is a sum of the following:

   1  Date
   2  Time
   4  Microsecond time
   8  Long filename
  16  Short filename

I.e. to prefix each log line with date and time, set -logflags=3 (1 + 2 from
above). The value 0 is used to disable all of the above. The default is to
show time only (2).


Development Settings
--------------------

The following environment variables modify syncthing's behavior in ways that
are mostly useful for developers. Use with care.

 STGUIASSETS     Directory to load GUI assets from. Overrides compiled in assets.

 STTRACE         A comma separated string of facilities to trace. The valid
                 facility strings are:

                 - "beacon"   (the beacon package)
                 - "discover" (the discover package)
                 - "events"   (the events package)
                 - "files"    (the files package)
                 - "http"     (the main package; HTTP requests)
                 - "locks"    (the sync package; trace long held locks)
                 - "net"      (the main package; connections & network messages)
                 - "model"    (the model package)
                 - "scanner"  (the scanner package)
                 - "stats"    (the stats package)
                 - "suture"   (the suture package; service management)
                 - "upnp"     (the upnp package)
                 - "xdr"      (the xdr package)
                 - "all"      (all of the above)

 STPROFILER      Set to a listen address such as "127.0.0.1:9090" to start the
                 profiler with HTTP access.

 STCPUPROFILE    Write a CPU profile to cpu-$pid.pprof on exit.

 STHEAPPROFILE   Write heap profiles to heap-$pid-$timestamp.pprof each time
                 heap usage increases.

 STBLOCKPROFILE  Write block profiles to block-$pid-$timestamp.pprof every 20
                 seconds.

 STPERFSTATS     Write running performance statistics to perf-$pid.csv. Not
                 supported on Windows.

 STNOUPGRADE     Disable automatic upgrades.

 GOMAXPROCS      Set the maximum number of CPU cores to use. Defaults to all
                 available CPU cores.

 GOGC            Percentage of heap growth at which to trigger GC. Default is
                 100. Lower numbers keep peak memory usage down, at the price
                 of CPU usage (ie. performance).`
)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:203

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:202
var (
	reset			bool
	showVersion		bool
	doUpgrade		bool
	doUpgradeCheck		bool
	upgradeTo		string
	noBrowser		bool
	noConsole		bool
	generateDir		string
	logFile			string
	auditEnabled		bool
	verbose			bool
	noRestart		= os.Getenv("STNORESTART") != "1"
	noUpgrade		= os.Getenv("STNOUPGRADE") != ""
	guiAddress		= os.Getenv("STGUIADDRESS")
	guiAuthentication	= os.Getenv("STGUIAUTH")
	guiAPIKey		= os.Getenv("STGUIAPIKEY")
	profiler		= os.Getenv("STPROFILER")
	guiAssets		= os.Getenv("STGUIASSETS")
	cpuProfile		= os.Getenv("STCPUPROFILE") != ""
	stRestarting		= os.Getenv("STRESTART") != ""
	innerProcess		= os.Getenv("STNORESTART") != "" || os.Getenv("STMONITORED") != ""
)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:227

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:226
func main() {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:229

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:228
	fmt.Printf("Say hello to Syncthing from Go!\n")
												_Cfunc_iosmain(0, nil)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:232

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:231
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:235

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:234
func WebServer() {
												fmt.Println("Run DevSync")
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:238

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:237
	if runtime.GOOS == "windows" {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:242

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:241
		flag.StringVar(&logFile, "logfile", "", "Log file name (use \"-\" for stdout)")
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:245

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:244
		flag.BoolVar(&noConsole, "no-console", false, "Hide console window")
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:248

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:247
	flag.StringVar(&generateDir, "generate", "", "Generate key and config in specified dir, then exit")
												flag.StringVar(&guiAddress, "gui-address", guiAddress, "Override GUI address")
												flag.StringVar(&guiAuthentication, "gui-authentication", guiAuthentication, "Override GUI authentication; username:password")
												flag.StringVar(&guiAPIKey, "gui-apikey", guiAPIKey, "Override GUI API key")
												flag.StringVar(&confDir, "home", "", "Set configuration directory")
												flag.IntVar(&logFlags, "logflags", logFlags, "Select information in log line prefix")
												flag.BoolVar(&noBrowser, "no-browser", false, "Do not start browser")
												flag.BoolVar(&noRestart, "no-restart", noRestart, "Do not restart; just exit")
												flag.BoolVar(&reset, "reset", false, "Reset the database")
												flag.BoolVar(&doUpgrade, "upgrade", false, "Perform upgrade")
												flag.BoolVar(&doUpgradeCheck, "upgrade-check", false, "Check for available upgrade")
												flag.BoolVar(&showVersion, "version", false, "Show version")
												flag.StringVar(&upgradeTo, "upgrade-to", upgradeTo, "Force upgrade directly from specified URL")
												flag.BoolVar(&auditEnabled, "audit", false, "Write events to audit file")
												flag.BoolVar(&verbose, "verbose", false, "Print verbose log output")
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:264

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:263
	flag.Usage = usageFor(flag.CommandLine, usage, fmt.Sprintf(extraUsage, baseDirs["config"]))
												flag.Parse()
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:267

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:266
	if noConsole {
		osutil.HideConsole()
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:271

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:270
	if confDir != "" {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:273

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:272
		baseDirs["config"] = confDir
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:276

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:275
	if err := expandLocations(); err != nil {
		l.Fatalln(err)
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:280

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:279
	if guiAssets == "" {
		guiAssets = locations[locGUIAssets]
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:284

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:283
	if runtime.GOOS == "windows" {
		if logFile == "" {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:287

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:286
			logFile = locations[locLogFile]
		} else if logFile == "-" {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:290

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:289
			logFile = ""
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:294

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:293
	if showVersion {
		fmt.Println(LongVersion)
		return
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:299

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:298
	l.SetFlags(logFlags)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:301

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:300
	if generateDir != "" {
		dir, err := osutil.ExpandTilde(generateDir)
		if err != nil {
			l.Fatalln("generate:", err)
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:307

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:306
		info, err := os.Stat(dir)
													if err == nil && !info.IsDir() {
			l.Fatalln(dir, "is not a directory")
		}
		if err != nil && os.IsNotExist(err) {
			err = osutil.MkdirAll(dir, 0700)
			if err != nil {
				l.Fatalln("generate:", err)
			}
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:318

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:317
		certFile, keyFile := filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem")
													cert, err := tls.LoadX509KeyPair(certFile, keyFile)
													if err == nil {
			l.Warnln("Key exists; will not overwrite.")
			l.Infoln("Device ID:", protocol.NewDeviceID(cert.Certificate[0]))
		} else {
			cert, err = newCertificate(certFile, keyFile, tlsDefaultCommonName)
			myID = protocol.NewDeviceID(cert.Certificate[0])
			if err != nil {
				l.Fatalln("load cert:", err)
			}
			if err == nil {
				l.Infoln("Device ID:", protocol.NewDeviceID(cert.Certificate[0]))
			}
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:334

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:333
		cfgFile := filepath.Join(dir, "config.xml")
													if _, err := os.Stat(cfgFile); err == nil {
			l.Warnln("Config exists; will not overwrite.")
			return
		}
		var myName, _ = os.Hostname()
		var newCfg = defaultConfig(myName)
		var cfg = config.Wrap(cfgFile, newCfg)
		err = cfg.Save()
		if err != nil {
			l.Warnln("Failed to save config", err)
		}
		return
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:351

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:350
	if info, err := os.Stat(baseDirs["config"]); err == nil && !info.IsDir() {
		l.Fatalln("Config directory", baseDirs["config"], "is not a directory")
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:356

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:355
	ensureDir(baseDirs["config"], 0700)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:358

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:357
	if upgradeTo != "" {
		err := upgrade.ToURL(upgradeTo)
		if err != nil {
			l.Fatalln("Upgrade:", err)
		}
		l.Okln("Upgraded from", upgradeTo)
		return
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:367

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:366
	if doUpgrade || doUpgradeCheck {
		rel, err := upgrade.LatestRelease(Version)
		if err != nil {
			l.Fatalln("Upgrade:", err)
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:373

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:372
		if upgrade.CompareVersions(rel.Tag, Version) <= 0 {
			l.Infof("No upgrade available (current %q >= latest %q).", Version, rel.Tag)
			os.Exit(exitNoUpgradeAvailable)
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:378

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:377
		l.Infof("Upgrade available (current %q < latest %q)", Version, rel.Tag)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:380

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:379
		if doUpgrade {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:382

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:381
			_, err = leveldb.OpenFile(locations[locDatabase], &opt.Options{OpenFilesCacheCapacity: 100})
														if err != nil {
				l.Infoln("Attempting upgrade through running Syncthing...")
				err = upgradeViaRest()
				if err != nil {
					l.Fatalln("Upgrade:", err)
				}
				l.Okln("Syncthing upgrading")
				return
			}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:393

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:392
			err = upgrade.To(rel)
														if err != nil {
				l.Fatalln("Upgrade:", err)
			}
			l.Okf("Upgraded to %q", rel.Tag)
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:400

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:399
		return
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:403

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:402
	if reset {
		resetDB()
		return
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:408

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:407
	if noRestart {
		syncthingMain()
	} else {
		monitorMain()
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:415

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:414
func upgradeViaRest() error {
	cfg, err := config.Load(locations[locConfigFile], protocol.LocalDeviceID)
	if err != nil {
		return err
	}
	target := cfg.GUI().Address
	if cfg.GUI().UseTLS {
		target = "https://" + target
	} else {
		target = "http://" + target
	}
												r, _ := http.NewRequest("POST", target+"/rest/system/upgrade", nil)
												r.Header.Set("X-API-Key", cfg.GUI().APIKey)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:429

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:428
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport:	tr,
		Timeout:	60 * time.Second,
	}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		bs, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		return errors.New(string(bs))
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:449

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:448
	return err
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:452

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:451
func syncthingMain() {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:455

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:454
	mainSvc := suture.New("main", suture.Spec{
		Log: func(line string) {
			if debugSuture {
				l.Debugln(line)
			}
		},
	})
												mainSvc.ServeBackground()
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:466

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:465
	l.SetPrefix("[start] ")
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:468

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:467
	if auditEnabled {
		startAuditing(mainSvc)
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:472

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:471
	if verbose {
		mainSvc.Add(newVerboseSvc())
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:477

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:476
	apiSub := events.NewBufferedSubscription(events.Default.Subscribe(events.AllEvents), 1000)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:479

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:478
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:484

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:483
	cert, err := tls.LoadX509KeyPair(locations[locCertFile], locations[locKeyFile])
												if err != nil {
		cert, err = newCertificate(locations[locCertFile], locations[locKeyFile], tlsDefaultCommonName)
		if err != nil {
			l.Fatalln("load cert:", err)
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:494

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:493
	predictableRandom.Seed(seedFromBytes(cert.Certificate[0]))
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:496

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:495
	myID = protocol.NewDeviceID(cert.Certificate[0])
												l.SetPrefix(fmt.Sprintf("[%s] ", myID.String()[:5]))
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:499

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:498
	l.Infoln(LongVersion)
												l.Infoln("My ID:", myID)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:504

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:503
	events.Default.Log(events.Starting, map[string]string{
		"home":	baseDirs["config"],
		"myID":	myID.String(),
	})
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:511

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:510
	cfgFile := locations[locConfigFile]
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:513

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:512
	var myName string
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:518

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:517
	if info, err := os.Stat(cfgFile); err == nil {
		if !info.Mode().IsRegular() {
			l.Fatalln("Config file is not a file?")
		}
		cfg, err = config.Load(cfgFile, myID)
		if err == nil {
			myCfg := cfg.Devices()[myID]
			if myCfg.Name == "" {
				myName, _ = os.Hostname()
			} else {
				myName = myCfg.Name
			}
		} else {
			l.Fatalln("Configuration:", err)
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:534

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:533
	} else {
		l.Infoln("No config file; starting with empty defaults")
		myName, _ = os.Hostname()
		newCfg := defaultConfig(myName)
		cfg = config.Wrap(cfgFile, newCfg)
		cfg.Save()
		l.Infof("Edit %s to taste or use the GUI\n", cfgFile)
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:543

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:542
	if cfg.Raw().OriginalVersion != config.CurrentVersion {
													l.Infoln("Archiving a copy of old config file format")
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:546

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:545
		osutil.Rename(cfgFile, cfgFile+fmt.Sprintf(".v%d", cfg.Raw().OriginalVersion))
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:548

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:547
		cfg.Save()
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:551

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:550
	if err := checkShortIDs(cfg); err != nil {
		l.Fatalln("Short device IDs are in conflict. Unlucky!\n  Regenerate the device ID of one if the following:\n  ", err)
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:555

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:554
	if len(profiler) > 0 {
		go func() {
			l.Debugln("Starting profiler on", profiler)
			runtime.SetBlockProfileRate(1)
			err := http.ListenAndServe(profiler, nil)
			if err != nil {
				l.Fatalln(err)
			}
		}()
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:569

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:568
	tlsCfg := &tls.Config{
		Certificates:		[]tls.Certificate{cert},
		NextProtos:		[]string{bepProtocolName},
		ClientAuth:		tls.RequestClientCert,
		SessionTicketsDisabled:	true,
		InsecureSkipVerify:	true,
		MinVersion:		tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		},
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:589

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:588
	opts := cfg.Options()
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:591

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:590
	if !opts.SymlinksEnabled {
		symlinks.Supported = false
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:595

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:594
	protocol.PingTimeout = time.Duration(opts.PingTimeoutS) * time.Second
												protocol.PingIdleTime = time.Duration(opts.PingIdleTimeS) * time.Second
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:598

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:597
	if opts.MaxSendKbps > 0 {
		writeRateLimit = ratelimit.NewBucketWithRate(float64(1000*opts.MaxSendKbps), int64(5*1000*opts.MaxSendKbps))
	}
	if opts.MaxRecvKbps > 0 {
		readRateLimit = ratelimit.NewBucketWithRate(float64(1000*opts.MaxRecvKbps), int64(5*1000*opts.MaxRecvKbps))
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:605

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:604
	if (opts.MaxRecvKbps > 0 || opts.MaxSendKbps > 0) && !opts.LimitBandwidthInLan {
		lans, _ = osutil.GetLans()
		networks := make([]string, 0, len(lans))
		for _, lan := range lans {
			networks = append(networks, lan.String())
		}
		l.Infoln("Local networks:", strings.Join(networks, ", "))
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:614

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:613
	dbFile := locations[locDatabase]
												ldb, err := leveldb.OpenFile(dbFile, dbOpts())
												if err != nil && errors.IsCorrupted(err) {
		ldb, err = leveldb.RecoverFile(dbFile, dbOpts())
	}
	if err != nil {
		l.Fatalln("Cannot open database:", err, "- Is another copy of Syncthing already running?")
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:624

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:623
	folders := cfg.Folders()
												for _, folder := range db.ListFolders(ldb) {
		if _, ok := folders[folder]; !ok {
			l.Infof("Cleaning data for dropped folder %q", folder)
			db.DropFolder(ldb, folder)
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:632

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:631
	m := model.NewModel(cfg, myID, myName, "syncthing", Version, ldb)
												cfg.Subscribe(m)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:635

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:634
	if t := os.Getenv("STDEADLOCKTIMEOUT"); len(t) > 0 {
		it, err := strconv.Atoi(t)
		if err == nil {
			m.StartDeadlockDetector(time.Duration(it) * time.Second)
		}
	} else if !IsRelease || IsBeta {
		m.StartDeadlockDetector(20 * 60 * time.Second)
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:647

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:646
	for _, folderCfg := range cfg.Folders() {
		m.AddFolder(folderCfg)
		for _, device := range folderCfg.DeviceIDs() {
			if device == myID {
				continue
			}
			m.Index(device, folderCfg.ID, nil, 0, nil)
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:657

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:656
		if folderCfg.ReadOnly {
			l.Okf("Ready to synchronize %s (read only; no external updates accepted)", folderCfg.ID)
			m.StartFolderRO(folderCfg.ID)
		} else {
			l.Okf("Ready to synchronize %s (read-write)", folderCfg.ID)
			m.StartFolderRW(folderCfg.ID)
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:666

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:665
	mainSvc.Add(m)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:670

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:669
	setupGUI(mainSvc, cfg, m, apiSub)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:674

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:673
	addr, err := net.ResolveTCPAddr("tcp", opts.ListenAddress[0])
												if err != nil {
		l.Fatalln("Bad listen address:", err)
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:681

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:680
	localPort := addr.Port
												discoverer = discovery(localPort)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:687

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:686
	if opts.UPnPEnabled {
		upnpSvc := newUPnPSvc(cfg, localPort)
		mainSvc.Add(upnpSvc)
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:692

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:691
	connectionSvc := newConnectionSvc(cfg, myID, m, tlsCfg)
												cfg.Subscribe(connectionSvc)
												mainSvc.Add(connectionSvc)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:696

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:695
	if cpuProfile {
		f, err := os.Create(fmt.Sprintf("cpu-%d.pprof", os.Getpid()))
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:705

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:704
	for _, device := range cfg.Devices() {
		if len(device.Name) > 0 {
			l.Infof("Device %s is %q at %v", device.DeviceID, device.Name, device.Addresses)
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:711

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:710
	if opts.URAccepted > 0 && opts.URAccepted < usageReportVersion {
		l.Infoln("Anonymous usage report has changed; revoking acceptance")
		opts.URAccepted = 0
		opts.URUniqueID = ""
		cfg.SetOptions(opts)
	}
	if opts.URAccepted >= usageReportVersion {
		if opts.URUniqueID == "" {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:721

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:720
			opts.URUniqueID = randomString(8)
														cfg.SetOptions(opts)
														cfg.Save()
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:730

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:729
	newUsageReportingManager(m, cfg)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:732

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:731
	if opts.RestartOnWakeup {
		go standbyMonitor()
												}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:736

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:735
	if opts.AutoUpgradeIntervalH > 0 {
		if noUpgrade {
			l.Infof("No automatic upgrades; STNOUPGRADE environment variable defined.")
		} else if IsRelease {
			go autoUpgrade()
		} else {
			l.Infof("No automatic upgrades; %s is not a release version.", Version)
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:746

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:745
	events.Default.Log(events.StartupComplete, map[string]string{
		"myID": myID.String(),
												})
												go generatePingEvents()
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:751

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:750
	cleanConfigDirectory()
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:753

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:752
	code := <-stop
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:755

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:754
	mainSvc.Stop()
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:757

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:756
	l.Okln("Exiting")
												os.Exit(code)
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:761

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:760
func dbOpts() *opt.Options {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:765

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:764
	blockCacheCapacity := 8 << 20
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:767

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:766
	const maxCapacity = 64 << 20
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:769

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:768
	const maxAtRAM = 8 << 30
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:771

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:770
	if v := cfg.Options().DatabaseBlockCacheMiB; v != 0 {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:773

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:772
		blockCacheCapacity = v << 20
	} else if bytes, err := memorySize(); err == nil {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:778

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:777
		if bytes > maxAtRAM {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:780

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:779
			blockCacheCapacity = maxCapacity
		} else if bytes > maxAtRAM/maxCapacity*int64(blockCacheCapacity) {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:783

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:782
			blockCacheCapacity = int(bytes * maxCapacity / maxAtRAM)
		}
		l.Infoln("Database block cache capacity", blockCacheCapacity/1024, "KiB")
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:788

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:787
	return &opt.Options{
		OpenFilesCacheCapacity:	100,
		BlockCacheCapacity:	blockCacheCapacity,
		WriteBuffer:		4 << 20,
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:795

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:794
func startAuditing(mainSvc *suture.Supervisor) {
	auditFile := timestampedLoc(locAuditLog)
	fd, err := os.OpenFile(auditFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		l.Fatalln("Audit:", err)
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:802

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:801
	auditSvc := newAuditSvc(fd)
												mainSvc.Add(auditSvc)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:807

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:806
	auditSvc.WaitForStart()
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:809

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:808
	l.Infoln("Audit log in", auditFile)
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:812

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:811
func setupGUI(mainSvc *suture.Supervisor, cfg *config.Wrapper, m *model.Model, apiSub *events.BufferedSubscription) {
												opts := cfg.Options()
												guiCfg := overrideGUIConfig(cfg.GUI(), guiAddress, guiAuthentication, guiAPIKey)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:816

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:815
	if guiCfg.Enabled && guiCfg.Address != "" {
		addr, err := net.ResolveTCPAddr("tcp", guiCfg.Address)
		if err != nil {
			l.Fatalf("Cannot start GUI on %q: %v", guiCfg.Address, err)
		} else {
			var hostOpen, hostShow string
			switch {
			case addr.IP == nil:
				hostOpen = "localhost"
				hostShow = "0.0.0.0"
			case addr.IP.IsUnspecified():
				hostOpen = "localhost"
				hostShow = addr.IP.String()
			default:
				hostOpen = addr.IP.String()
				hostShow = hostOpen
			}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:834

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:833
			var proto = "http"
														if guiCfg.UseTLS {
				proto = "https"
			}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:839

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:838
			urlShow := fmt.Sprintf("%s://%s/", proto, net.JoinHostPort(hostShow, strconv.Itoa(addr.Port)))
														l.Infoln("Starting web GUI on", urlShow)
														api, err := newAPISvc(myID, guiCfg, guiAssets, m, apiSub)
														if err != nil {
				l.Fatalln("Cannot start GUI:", err)
			}
														cfg.Subscribe(api)
														mainSvc.Add(api)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:848

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:847
			if opts.StartBrowser && !noBrowser && !stRestarting {
															urlOpen := fmt.Sprintf("%s://%s/", proto, net.JoinHostPort(hostOpen, strconv.Itoa(addr.Port)))
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:852

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:851
				go openURL(urlOpen)
			}
		}
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:858

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:857
func defaultConfig(myName string) config.Configuration {
	newCfg := config.New(myID)
	newCfg.Folders = []config.FolderConfiguration{
		{
			ID:			"default",
			RawPath:		locations[locDefFolder],
			RescanIntervalS:	60,
			Devices:		[]config.FolderDeviceConfiguration{{DeviceID: myID}},
		},
	}
	newCfg.Devices = []config.DeviceConfiguration{
		{
			DeviceID:	myID,
			Addresses:	[]string{"dynamic"},
			Name:		myName,
		},
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:876

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:875
	port, err := getFreePort("127.0.0.1", 8384)
												if err != nil {
		l.Fatalln("get free port (GUI):", err)
	}
												newCfg.GUI.Address = fmt.Sprintf("127.0.0.1:%d", port)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:882

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:881
	port, err = getFreePort("0.0.0.0", 22000)
												if err != nil {
		l.Fatalln("get free port (BEP):", err)
	}
	newCfg.Options.ListenAddress = []string{fmt.Sprintf("0.0.0.0:%d", port)}
	return newCfg
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:890

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:889
func generatePingEvents() {
	for {
		time.Sleep(pingEventInterval)
		events.Default.Log(events.Ping, nil)
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:897

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:896
func resetDB() error {
	return os.RemoveAll(locations[locDatabase])
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:903

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:902
func Restart() {
	l.Infoln("Restarting")
	stop <- exitRestarting
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:908

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:907
func restart() {
	l.Infoln("Restarting")
	stop <- exitRestarting
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:914

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:913
func ShutDown() {
	l.Infoln("Shutting down")
	stop <- exitSuccess
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:920

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:919
func shutdown() {
	l.Infoln("Shutting down")
	stop <- exitSuccess
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:925

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:924
func discovery(extPort int) *discover.Discoverer {
												opts := cfg.Options()
												disc := discover.NewDiscoverer(myID, opts.ListenAddress)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:929

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:928
	if opts.LocalAnnEnabled {
		l.Infoln("Starting local discovery announcements")
		disc.StartLocal(opts.LocalAnnPort, opts.LocalAnnMCAddr)
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:934

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:933
	if opts.GlobalAnnEnabled {
		l.Infoln("Starting global discovery announcements")
		disc.StartGlobal(opts.GlobalAnnServers, uint16(extPort))
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:939

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:938
	return disc
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:942

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:941
func ensureDir(dir string, mode int) {
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := osutil.MkdirAll(dir, 0700)
		if err != nil {
			l.Fatalln(err)
		}
	} else if mode >= 0 && err == nil && int(fi.Mode()&0777) != mode {
													err := os.Chmod(dir, os.FileMode(mode))
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:952

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:951
		if err != nil {
			l.Warnln(err)
		}
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:961

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:960
func getFreePort(host string, ports ...int) (int, error) {
	for _, port := range ports {
		c, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			c.Close()
			return port, nil
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:970

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:969
	c, err := net.Listen("tcp", host+":0")
												if err != nil {
		return 0, err
	}
	addr := c.Addr().(*net.TCPAddr)
	c.Close()
	return addr.Port, nil
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:979

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:978
func overrideGUIConfig(cfg config.GUIConfiguration, address, authentication, apikey string) config.GUIConfiguration {
	if address != "" {
													cfg.Enabled = true
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:983

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:982
		if !strings.Contains(address, "//") {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:985

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:984
			cfg.Address = address
		} else {
			parsed, err := url.Parse(address)
			if err != nil {
				l.Fatalln(err)
			}
			cfg.Address = parsed.Host
			switch parsed.Scheme {
			case "http":
				cfg.UseTLS = false
			case "https":
				cfg.UseTLS = true
			default:
				l.Fatalln("Unknown scheme:", parsed.Scheme)
			}
		}
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1003

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1002
	if authentication != "" {
													authenticationParts := strings.SplitN(authentication, ":", 2)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1006

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1005
		hash, err := bcrypt.GenerateFromPassword([]byte(authenticationParts[1]), 0)
													if err != nil {
			l.Fatalln("Invalid GUI password:", err)
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1011

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1010
		cfg.User = authenticationParts[0]
													cfg.Password = string(hash)
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1015

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1014
	if apikey != "" {
		cfg.APIKey = apikey
	}
	return cfg
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1021

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1020
func standbyMonitor() {
	restartDelay := time.Duration(60 * time.Second)
	now := time.Now()
	for {
		time.Sleep(10 * time.Second)
		if time.Since(now) > 2*time.Minute {
														l.Infof("Paused state detected, possibly woke up from standby. Restarting in %v.", restartDelay)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1032

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1031
			time.Sleep(restartDelay)
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1034

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1033
			restart()
														return
		}
		now = time.Now()
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1041

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1040
func autoUpgrade() {
	timer := time.NewTimer(0)
	sub := events.Default.Subscribe(events.DeviceConnected)
	for {
		select {
		case event := <-sub.C():
			data, ok := event.Data.(map[string]string)
			if !ok || data["clientName"] != "syncthing" || upgrade.CompareVersions(data["clientVersion"], Version) != upgrade.Newer {
				continue
			}
			l.Infof("Connected to device %s with a newer version (current %q < remote %q). Checking for upgrades.", data["id"], Version, data["clientVersion"])
		case <-timer.C:
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1055

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1054
		rel, err := upgrade.LatestRelease(Version)
													if err == upgrade.ErrUpgradeUnsupported {
			events.Default.Unsubscribe(sub)
			return
		}
		if err != nil {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1063

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1062
			l.Infoln("Automatic upgrade:", err)
														timer.Reset(time.Duration(cfg.Options().AutoUpgradeIntervalH) * time.Hour)
														continue
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1068

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1067
		if upgrade.CompareVersions(rel.Tag, Version) != upgrade.Newer {
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1070

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1069
			timer.Reset(time.Duration(cfg.Options().AutoUpgradeIntervalH) * time.Hour)
														continue
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1074

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1073
		l.Infof("Automatic upgrade (current %q < latest %q)", Version, rel.Tag)
													err = upgrade.To(rel)
													if err != nil {
			l.Warnln("Automatic upgrade:", err)
			timer.Reset(time.Duration(cfg.Options().AutoUpgradeIntervalH) * time.Hour)
			continue
		}
		events.Default.Unsubscribe(sub)
		l.Warnf("Automatically upgraded to version %q. Restarting in 1 minute.", rel.Tag)
		time.Sleep(time.Minute)
		stop <- exitUpgrading
		return
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1091

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1090
func cleanConfigDirectory() {
	patterns := map[string]time.Duration{
		"panic-*.log":		7 * 24 * time.Hour,
		"audit-*.log":		7 * 24 * time.Hour,
		"index":		14 * 24 * time.Hour,
		"config.xml.v*":	30 * 24 * time.Hour,
		"*.idx.gz":		30 * 24 * time.Hour,
		"backup-of-v0.8":	30 * 24 * time.Hour,
	}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1101

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1100
	for pat, dur := range patterns {
		pat = filepath.Join(baseDirs["config"], pat)
		files, err := osutil.Glob(pat)
		if err != nil {
			l.Infoln("Cleaning:", err)
			continue
		}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1109

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1108
		for _, file := range files {
			info, err := osutil.Lstat(file)
			if err != nil {
				l.Infoln("Cleaning:", err)
				continue
			}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1116

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1115
			if time.Since(info.ModTime()) > dur {
				if err = os.RemoveAll(file); err != nil {
					l.Infoln("Cleaning:", err)
				} else {
					l.Infoln("Cleaned away old file", filepath.Base(file))
				}
			}
		}
	}
}
//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1130

//line /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/cmd/syncthing/main.go:1129
func checkShortIDs(cfg *config.Wrapper) error {
	exists := make(map[uint64]protocol.DeviceID)
	for deviceID := range cfg.Devices() {
		shortID := deviceID.Short()
		if otherID, ok := exists[shortID]; ok {
			return fmt.Errorf("%v in conflict with %v", deviceID, otherID)
		}
		exists[shortID] = deviceID
	}
	return nil
}
