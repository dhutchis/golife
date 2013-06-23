// glife project main.go
package main

import (
	"fmt"
	"github.com/denine99/glife/glife"
)

func main() {
	fmt.Println("Hello World!")
	x := glife.CreateAllDead(2, 3)
	fmt.Printf("%T: %v\n", x, x)
}
