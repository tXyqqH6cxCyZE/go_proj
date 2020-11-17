package main

import (
	"fmt"
	"unsafe"
)

/**
返回局部变量指针是安全的，编译器会根据需要将其分配在GC Heap上
*/
func test() *int {
	x := 100
	return &x // 使用runtime.new 分配 x 内存。但在内联时，也可能直接分配在目标栈。
}

/**
将pointer 转换成 uintptr，可变相实现指针运算。
1. unsafe.Pointer() : 通用指针类型，用于转换不同类型的指针，不能进行指针运算
2. uintptr : 用于指针运算，GC不把uintptr当指针，uintptr无法持有对象。uintptr类型的目标会被回收
*/

func main() {
	d := struct {
		s string
		x int
		l int
	}{"abc", 100, 300}

	p := uintptr(unsafe.Pointer(&d)) // *struct -> Pointer -> uintptr
	p += unsafe.Offsetof(d.l)        // uintptr + offset

	p2 := unsafe.Pointer(p) // uintptr -> Pointer
	px := (*int)(p2)        // Pointer -> *int
	*px = 400

	fmt.Printf("%#v\n", d)
}
