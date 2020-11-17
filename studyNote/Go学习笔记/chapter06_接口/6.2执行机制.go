package main

import "fmt"

/**
接口对象由接口表（interface table）指针和数据指针组成。
runtime.h
struct Iface
{
	Itab* tab;
	void* data;
}

struct Itab
{
	InterfaceType* inter;
	Type*		   type;
	void (*fun[])(void);
}

接口表：存储元数据信息，包括接口类型，动态类型，以及实现接口的方法指针。
无论是反射还是通过接口调用方法，都会用到这些信息。

数据指针：持有的是目标对象的只读复制品，复制完整对象或指针

补充：只有tab 和 data都为nil时，接口才等于niL。
*/

type User62 struct {
	id   int
	name string
}

func main() {
	// 数据指针的使用
	u := User62{1, "Tom"}
	var i interface{} = u

	u.id = 2
	u.name = "Jack"

	fmt.Printf("%v\n", u)
	fmt.Printf("%v\n", i.(User62))

	// 接口转型返回临时对象，只有使用指针才能修改其状态。
	u1 := User62{1, "Tom"}
	var vi, pi interface{} = u1, &u1

	// vi.(User).name = "Jack"		// Error: cannot assign to vi.(User).name
	pi.(*User62).name = "Jack"

	fmt.Printf("%v\n", vi.(User62))
	fmt.Printf("%v\n", pi.(*User62))

}
