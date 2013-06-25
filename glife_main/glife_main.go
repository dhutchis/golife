// glife project main.go
package main

import (
	"fmt"
	"github.com/denine99/glife/glife"
	"os"
)

func main() {
	fmt.Println("Hello World!")
	x := glife.CreateAllDead(2, 3)
	fmt.Printf("%T:\n%v\n", x, x)

	var f glife.Field
	file, err := os.Open("../cells/beacon.cells")
	if err != nil {
		fmt.Println("can't find beacon.cells file")
		return
	}
	defer file.Close()
	f = glife.ReadFieldFrom(file)
	fmt.Printf("%T:\n%v\n", f, f)

	f.Run(1)
	fmt.Println(f)
	f.Run(1)
	fmt.Println(f)
	f.Run(1)
	fmt.Println(f)
	f.Run(1)
	fmt.Println(f)
}
