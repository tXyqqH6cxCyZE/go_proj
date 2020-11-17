package main

import "fmt"

/**
可以像字段成员那样访问匿名字段方法，编译器负责查找
*/
type User struct {
	id   int
	name string
}

func (self *User) ToString() string { // receiver = &{Manager.User}
	return fmt.Sprintf("User: %p, %v", self, self)
}

type Manager struct {
	User
}

/**
通过匿名字段，可以获得和继承类似的复用能力。
依据编译器查找次序，只需在外层定义同名方法，就可以实现“override”
*/
type User1 struct {
	id   int
	name string
}

type Manager1 struct {
	User
	title string
}

func (self *User1) ToString() string {
	return fmt.Sprintf("User1: %p, %v", self, self)
}

func (self *Manager1) ToString() string {
	return fmt.Sprintf("Manager1: %p, %v", self, self)
}

func main() {
	m := Manager{User{1, "Tom"}}
	fmt.Printf("Manager: %p\n", &m)
	fmt.Println(m.ToString())

	m1 := Manager1{User{1, "Tom"}, "Administrator"}

	fmt.Println(m1.ToString())
	fmt.Println(m1.User.ToString())
}
