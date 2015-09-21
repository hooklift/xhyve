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
