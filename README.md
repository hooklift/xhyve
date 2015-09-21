# libxhyve
Go bindings to use xhyve as a library.

### Install
go get github.com/hooklift/libxhyve

### Example

```go
package main

import (
	"os"
	"github.com/hooklift/libxhyve"
)

func main() {
	if err := xhyve.Run(os.Args); err != nil {
		panic(err)
	}
}
```

```bash
sudo ./nnn -m 1024M -c 1 -A -s 0:0,hostbridge -s 31,lpc \
-l com1,stdio -s 2:0,virtio-net -U 6BCE442E-4359-4BD9-84F7-EDFB8EC6D2EF \
-f 'kexec,imgs/stable.766.3.0.coreos_production_pxe.vmlinuz,imgs/stable.766.3.0.coreos_production_pxe_image.cpio.gz,earlyprintk=serial console=ttyS0 coreos.autologin'
```
