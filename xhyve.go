// +build darwin

package main

// #cgo CFLAGS: -I${SRCDIR}/vendor/xhyve/include -x c -std=c11 -fno-common -arch x86_64 -DXHYVE_CONFIG_ASSERT -lxhyve -Os -fstrict-aliasing -Wno-unknown-warning-option -Wno-reserved-id-macro -pedantic -fmessage-length=152 -fdiagnostics-show-note-include-stack -fmacro-backtrace-limit=0
// #cgo LDFLAGS: -L${SRCDIR} -lxhyve -arch x86_64 -framework Hypervisor -framework vmnet -force_load ${SRCDIR}/libxhyve.a
// #include <xhyve/xhyve.h>
// #include <string.h>
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

func init() {
	runtime.LockOSThread()
	//signal.Ignore()
}

// RunXHyve runs xhyve hypervisor.
func RunXHyve(params []string) error {
	argc := C.int(len(params))
	argv := make([]*C.char, argc)
	for i, arg := range params {
		argv[i] = C.CString(arg)
	}

	if err := C.run_xhyve(argc, &argv[0]); err != 0 {
		fmt.Printf("ERROR => %s\n", C.GoString(C.strerror(err)))
		return fmt.Errorf("Error initializing hypervisor")
	}

	//FIXME(c4milo): This is never reached right now due to xhyve calling
	// exit()
	for _, arg := range argv {
		C.free(unsafe.Pointer(arg))
	}

	return nil
}
