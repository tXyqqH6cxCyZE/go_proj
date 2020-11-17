package main

import "fmt"

/**
根据调用者不同，方法分为两种表现形式：
	instance.method(args...) ---> <type>.func(instance, args...)
	前者称为method value，后者method expression

两者都可像普通函数那样赋值和传参,
区别在于method value绑定实例，而method expression则须显式传参。
*/
type User4 struct {
	id   int
	name string
}

func (self *User4) Test() {
	fmt.Printf("%p, %v\n", self, self)
}

// 需要注意，method value 会复制receiver
type User41 struct {
	id   int
	name string
}

func (self User41) Test() {
	fmt.Println(self)
}

func main() {
	u := User4{1, "Tom"}
	u.Test()

	mValue := u.Test
	mValue() // 隐式传递 receiver

	mExpression := (*User4).Test
	mExpression(&u) // 显式传递 receiver

	fmt.Println("----------------------------------------------------")

	u1 := User41{1, "Tom"}
	mValue1 := u1.Test // 立即复制receiver,因为不是指针类型，不受后续修改影响
	// 在汇编层面，method value和闭包的实现方式相同，实际返回FuncVal类型对象。
	u1.id, u1.name = 2, "Jack"
	u1.Test()
	mValue1()
}
