/* Created by cgo - DO NOT EDIT. */
#include "_cgo_export.h"

extern void crosscall2(void (*fn)(void *, int), void *, int);
extern void _cgo_wait_runtime_init_done();

extern void _cgoexp_34a055b93e51_WebServer(void *, int);

void WebServer()
{
	_cgo_wait_runtime_init_done();
	struct {
		char unused;
	} __attribute__((__packed__)) a;
	crosscall2(_cgoexp_34a055b93e51_WebServer, &a, 0);
}
extern void _cgoexp_34a055b93e51_Restart(void *, int);

void Restart()
{
	_cgo_wait_runtime_init_done();
	struct {
		char unused;
	} __attribute__((__packed__)) a;
	crosscall2(_cgoexp_34a055b93e51_Restart, &a, 0);
}
extern void _cgoexp_34a055b93e51_ShutDown(void *, int);

void ShutDown()
{
	_cgo_wait_runtime_init_done();
	struct {
		char unused;
	} __attribute__((__packed__)) a;
	crosscall2(_cgoexp_34a055b93e51_ShutDown, &a, 0);
}
