package main

import (
	"fmt"
)

/**
Struct
- 值类型，赋值和传值会复制全部内容。可用“_”定义补位字段，支持指向自身类型的指针成员。
- 顺序初始化必须包含全部字段，否则会出错
- 支持匿名结构，可用作结构成员或定义变量
- 支持“==”，“!=”相关操作符，可用作map键类型
*/

type Node struct {
	_    int
	id   int
	data *byte
	next *Node
}

func main() {
	n1 := Node{
		id:   1,
		data: nil,
	}

	n2 := Node{
		id:   2,
		data: nil,
		next: &n1,
	}
	fmt.Print(n2)

	// 空结构“节省”内存，比如用来实现set数据结构，或者实现没有“状态”只有方法的"静态类"
	var null struct{}

	set := make(map[string]struct{})
	set["a"] = null

	/**
	1 匿名字段
	- 匿名字段不过是一种语法糖，从根本上说，就是一个与成员类型同名(不含包名)的字段。
	  被匿名嵌入的可以是任意类型，当然也包括指针。
	-
	*/
	type User struct {
		name string
	}

	type Manager struct {
		User
		title string
	}

	m := Manager{
		User:  User{"Tom"}, // 匿名字段的显式字段名，和类型名相同
		title: "Administrator",
	}
	fmt.Print(m)

	// 可以像普通字段那样访问匿名字段成员，编译器从外向内逐级查找所有层次的匿名字段，直到发现目标或出错
	type Resource struct {
		id int
	}

	type User1 struct {
		Resource
		name string
	}

	type Manager1 struct {
		User1
		title string
	}

	var m1 Manager1
	m1.id = 1
	m1.name = "Jack"
	m1.title = "Administator"

	// 外层同名字段会遮蔽嵌入字段成员，相同层次的同名字段也会让编译器无所适从。解决方案是使用显式字段名。
	type Resource2 struct {
		id   int
		name string
	}

	type Classify struct {
		id int
	}

	type User2 struct {
		Resource2        // Resource.id 与 Classify.id 处于同一层次
		Classify         //
		name      string // 遮蔽 Resource.name
	}

	u2 := User2{
		Resource2{1, "people"},
		Classify{100},
		"Jack",
	}

	println(u2.name)           // User.name: Jack
	println(u2.Resource2.name) // people
	//println(u2.id)					// Error: ambiguous selector u.id
	println(u2.Classify.id) // 100

	// 不能同时嵌入某一类型和其指针类型，因为它们名字相同
	type Resource3 struct {
		id int
	}

	type User3 struct {
		*Resource3
		// Resource3		// Error: duplicate field Resource
		name string
	}

	u3 := User3{
		&Resource3{1},
		"Administrator",
	}

	println(u3.id)
	println(u3.Resource3.id)

	/**
		2 面向对象
		- 面向对象三大特征里，Go仅支持封装，尽管匿名字段里的内存布局和行为类似继承。
	      没有class关键字，没有继承，多态等等
	*/
	type User4 struct {
		id   int
		name string
	}

	type Manager4 struct {
		User4
		title string
	}

	m4 := Manager4{
		User4{1, "Tom"},
		"Administrator",
	}
	//var u4 User4 = m4				Error: cannot use m(type Manager) as type User in assignment
	// 没有继承，自然也不会有多态
	var u4 User4 = m4.User4 // 同类型拷贝
	fmt.Println(u4)
}
