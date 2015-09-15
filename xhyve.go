package main

// #cgo CFLAGS: -I${SRCDIR}/vendor/xhyve/include
// #cgo LDFLAGS: -L${SRCDIR} -lxhyve -arch x86_64 -framework Hypervisor -framework vmnet
// #include "xhyve.h"
import "C"

func runXHyve() {

}
