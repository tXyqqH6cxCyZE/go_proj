package main

import (
	"fmt"
	"unsafe"
)

const (
	a = "1"
	b = len(a)
	c = unsafe.Sizeof(b)
)

/**
为什么是反过来的？
*/
func main() {
	x := 0x12345678

	p := unsafe.Pointer(&x)
	n := (*[4]byte)(p)

	for i := 0; i < len(n); i++ {
		fmt.Printf("%X ", n[i])
	}
}
