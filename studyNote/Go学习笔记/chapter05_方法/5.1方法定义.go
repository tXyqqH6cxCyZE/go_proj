package main

import "fmt"

/**
方法定义：
	方法总是绑定对象实例，并隐式将实例作为第一实参（receiver）
	- 只能为当前包内命名类型定义方法
	- 参数receiver可任意命名。如方法中未曾使用，可省略参数名
	- 参数receiver类型可以是T或*T。基类型不能是接口或指针
	- 不支持方法重载，receiver只是参数签名的组成部分
	- 可用实例value或pointer调用全部方法，编译器自动转换。
没有构造和析构方法，通常用简单工厂模式返回对象实例
*/

type Queue struct {
	elements []interface{}
}

func NewQueue() *Queue { // 创建对象实例
	return &Queue{make([]interface{}, 10)}
}

func (*Queue) Push(e interface{}) error { // 省略 receiver 参数名
	panic("not implemented")
}

//func (Queue) Push(e int) error {			// method redecleared: Queue.Push
//
//}

func (self *Queue) length() int { // receiver 参数名可以是 self, this或其他
	return len(self.elements)
}

/**
方法只不过是一种特殊的函数，只需将其还原，就知道receiver T 和 *T的差别
*/

type Data struct {
	x int
}

func (self Data) ValueTest() { // func ValueTest(self Data);
	fmt.Printf("Value: %p\n", &self)
}

func (self *Data) PointerTest() { // func PointerTest(self *Data)
	fmt.Printf("Pointer: %p\n", self)
}

func main() {
	d := Data{} // 这里返回的是结构体的值
	p := &d
	fmt.Printf("Data: %p\n", p)

	d.ValueTest()   // ValueTest(d)
	d.PointerTest() // PointerTest(&d)

	p.ValueTest()   // ValueTest(*p)
	p.PointerTest() // PointerTest(p)
}
