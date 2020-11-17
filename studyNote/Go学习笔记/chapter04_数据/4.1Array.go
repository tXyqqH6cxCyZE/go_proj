package main

import "fmt"

/**
1. 数组是值类型，赋值和传参会复制整个数组，而不是指针
2. 数组长度必须是常量，且是类型的组成部分。[2]int 和 [3]int是不同类型
3. 支持"==","!="操作符，因为内存总是被初始化过
4. 指针数组[n]*T，数组指针 *[n]T.
*/

var a = [3]int{1, 2}           // 未初始化元素值为0.
var b = [...]int{1, 2, 3, 4}   // 通过初始化值确定数组长度
var c = [5]int{2: 100, 4: 200} // 通过索引号初始化元素

var d = [...]struct { // 声明结构体数组，并初始化
	name string
	age  uint8
}{
	{"user1", 10}, // 可省略元素类型
	{"user2", 20}, // 别忘了最后一行的逗号。
}

//（因为特性1，值拷贝行为会造成性能问题，通常会建议使用slice，或数组指针）
func test(x [2]int) {
	fmt.Printf("x : %p\n", &x)
	x[1] = 1000
}

func main() {
	a := [2]int{}
	fmt.Printf("a : %p\n", &a)
	test(a)
	fmt.Println(a) // 可以看到两者地址不同，不是同一片内存
}
