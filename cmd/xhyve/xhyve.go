package main

import (
	"fmt"
	"os"

	"github.com/hooklift/xhyve"
)

func main() {
	if err := xhyve.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
