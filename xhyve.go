// +build darwin

package xhyve

// #cgo CFLAGS: -I${SRCDIR}/vendor/xhyve/include -x c -std=c11 -fno-common -arch x86_64 -DXHYVE_CONFIG_ASSERT -DVERSION=v0.2.0 -Os -fstrict-aliasing -Wno-unknown-warning-option -Wno-reserved-id-macro -pedantic -fmessage-length=152 -fdiagnostics-show-note-include-stack -fmacro-backtrace-limit=0
// #cgo LDFLAGS: -L${SRCDIR} -arch x86_64 -framework Hypervisor -framework vmnet
// #include <xhyve/xhyve.h>
// #include <string.h>
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

var argv []*C.char

// go_callback_exit gets called from within xhyve.c whenever a trap
// suspending the VM is trigger. This is so we can clean up resources
// in Go land.
//export go_callback_exit
func go_callback_exit(status C.int) {
	fmt.Printf("Exiting with error code %d\n", status)
	fmt.Printf("Releasing allocated memory from Go land... ")
	for _, arg := range argv {
		C.free(unsafe.Pointer(arg))
	}
	fmt.Println("done")

	C.exit_mevent_dispatch_loop = true
	runtime.UnlockOSThread()
}

func init() {
	runtime.LockOSThread()
}

// Run runs xhyve hypervisor.
func Run(params []string) error {
	argc := C.int(len(params))
	argv = make([]*C.char, argc)
	for i, arg := range params {
		argv[i] = C.CString(arg)
	}

	if err := C.run_xhyve(argc, &argv[0]); err != 0 {
		fmt.Printf("ERROR => %s\n", C.GoString(C.strerror(err)))
		return fmt.Errorf("Error initializing hypervisor")
	}

	return nil
}
