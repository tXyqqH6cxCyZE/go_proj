package main

import "fmt"

/**
接口转换：
	- 利用类型推断，可判断接口对象是否某个具体的接口或类型。
	- 还可用switch做批量类型判断，不支持fallthrough。
	- 超集接口对象可转换为子集接口，反之出错。
*/

type User63 struct {
	id   int
	name string
}

func (self *User63) String() string {
	return fmt.Sprintf("%d, %s", self.id, self.name)
}

func main() {
	var o interface{} = &User63{1, "Tom"}

	if i, ok := o.(fmt.Stringer); ok { // ok-idiom
		fmt.Println(i)
	}

	u := o.(*User63)
	// u := o.(User)		panic: interface is *match.User, not studyNote.User
	fmt.Println(u)
}
