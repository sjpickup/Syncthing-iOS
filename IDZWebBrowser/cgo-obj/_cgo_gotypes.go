// Created by cgo - DO NOT EDIT

package main

import "unsafe"

import _ "runtime/cgo"

import "syscall"

var _ syscall.Errno
func _Cgo_ptr(ptr unsafe.Pointer) unsafe.Pointer { return ptr }

//go:linkname _Cgo_always_false runtime.cgoAlwaysFalse
var _Cgo_always_false bool
//go:linkname _Cgo_use runtime.cgoUse
func _Cgo_use(interface{})
type _Ctype_char int8

type _Ctype_int int32

type _Ctype_void [0]byte

//go:linkname _cgo_runtime_cgocall runtime.cgocall
func _cgo_runtime_cgocall(unsafe.Pointer, uintptr) int32

//go:linkname _cgo_runtime_cmalloc runtime.cmalloc
func _cgo_runtime_cmalloc(uintptr) unsafe.Pointer

//go:linkname _cgo_runtime_cgocallback runtime.cgocallback
func _cgo_runtime_cgocallback(unsafe.Pointer, unsafe.Pointer, uintptr)

//go:cgo_import_static _cgo_34a055b93e51_Cfunc_iosmain
//go:linkname __cgofn__cgo_34a055b93e51_Cfunc_iosmain _cgo_34a055b93e51_Cfunc_iosmain
var __cgofn__cgo_34a055b93e51_Cfunc_iosmain byte
var _cgo_34a055b93e51_Cfunc_iosmain = unsafe.Pointer(&__cgofn__cgo_34a055b93e51_Cfunc_iosmain)

func _Cfunc_iosmain(p0 _Ctype_int, p1 **_Ctype_char) (r1 _Ctype_void) {
	_cgo_runtime_cgocall(_cgo_34a055b93e51_Cfunc_iosmain, uintptr(unsafe.Pointer(&p0)))
	if _Cgo_always_false {
		_Cgo_use(p0)
		_Cgo_use(p1)
	}
	return
}
//go:cgo_export_dynamic WebServer
//go:linkname _cgoexp_34a055b93e51_WebServer _cgoexp_34a055b93e51_WebServer
//go:cgo_export_static _cgoexp_34a055b93e51_WebServer
//go:nosplit
func _cgoexp_34a055b93e51_WebServer(a unsafe.Pointer, n int32) {	fn := WebServer
	_cgo_runtime_cgocallback(**(**unsafe.Pointer)(unsafe.Pointer(&fn)), a, uintptr(n));
}
//go:cgo_export_dynamic Restart
//go:linkname _cgoexp_34a055b93e51_Restart _cgoexp_34a055b93e51_Restart
//go:cgo_export_static _cgoexp_34a055b93e51_Restart
//go:nosplit
func _cgoexp_34a055b93e51_Restart(a unsafe.Pointer, n int32) {	fn := Restart
	_cgo_runtime_cgocallback(**(**unsafe.Pointer)(unsafe.Pointer(&fn)), a, uintptr(n));
}
//go:cgo_export_dynamic ShutDown
//go:linkname _cgoexp_34a055b93e51_ShutDown _cgoexp_34a055b93e51_ShutDown
//go:cgo_export_static _cgoexp_34a055b93e51_ShutDown
//go:nosplit
func _cgoexp_34a055b93e51_ShutDown(a unsafe.Pointer, n int32) {	fn := ShutDown
	_cgo_runtime_cgocallback(**(**unsafe.Pointer)(unsafe.Pointer(&fn)), a, uintptr(n));
}
