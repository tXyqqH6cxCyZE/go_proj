package main

import "fmt"

/**
接口定义：
	接口是一个或多个方法签名的集合，任何类型的方法集中只要拥有与之对应的全部方法，
	就表示它“实现”了该接口，无须在该类型上显式添加接口声明。

	所谓对应方法，是指有相同名称，参数列表（不包括参数名）以及返回值。
	当然，该类型还可以有其他方法。

	- 接口命名习惯以er结尾，结构体
	- 接口只有方法签名，没有实现
	- 接口没有数据字段
	- 可在接口中嵌入其他接口
	- 类型可实现多个接口

	补充：空接口interface{}没有任何方法签名，也就意味着任何类型都实现了空接口。
		 其作用类似面向对象语言中的根对象Object
*/

type Stringer interface {
	String() string
}

type Printer interface {
	Stringer // 接口嵌入
	Print()
}

type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("user %d, %s", self.id, self.name)
}

func (self *User) Print() {
	fmt.Println(self.String())
}

// 匿名接口可用作变量类型，或结构成员
type Tester struct {
	s interface {
		String() string
	}
}

type User6 struct {
	id   int
	name string
}

func (self *User6) String() string {
	return fmt.Sprintf("user %d %s", self.id, self.name)
}

func main() {
	var t Printer = &User{1, "Tom"} // *User 方法集包含String，Print。
	t.Print()

	t1 := Tester{&User6{1, "void"}}
	fmt.Println(t1.s.String())
}
