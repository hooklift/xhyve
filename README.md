# libxhyve
Go bindings to use xhyve as a library.

### Install
go get github.com/hooklift/xhyve

### Example

```go
package main

import (
	"os"
	"github.com/hooklift/xhyve"
)

func main() {
	if err := xhyve.Run(os.Args); err != nil {
		panic(err)
	}
}
```

There is small CLI that you can use to test the library.

```bash
cd cmd/xhyve; go build
sudo ./xhyve -m 1024M -c 1 -A -s 0:0,hostbridge -s 31,lpc \
-l com1,stdio -s 2:0,virtio-net -U 6BCE442E-4359-4BD9-84F7-EDFB8EC6D2EF \
-f 'kexec,imgs/vmlinuz,imgs/initrd.gz,earlyprintk=serial console=ttyS0'
```

![](https://asciinema.org/a/31xw6hl4amx8u5v453o20twf6)
