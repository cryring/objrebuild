package main

import (
	"fmt"
	"os"

	"github.com/cryring/objrebuild/obj"
)

func main() {
	argc := len(os.Args)
	if argc != 3 {
		fmt.Println("usage : objrebuild input.obj output.obj")
		return
	}

	iObj := os.Args[1]
	oObj := os.Args[2]

	o := obj.NewObj()
	if err := o.Load(iObj); err != nil {
		fmt.Println(err)
		return
	}
	if err := o.Save(oObj); err != nil {
		fmt.Println(err)
		return
	}
}
