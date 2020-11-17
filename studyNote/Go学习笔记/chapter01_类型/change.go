package main

import "fmt"

type bigint int64
type mysilce []int

func main() {
	x := 1234
	var b bigint = bigint(x) // 必须显式转换
	var b2 int64 = int64(b)
	fmt.Println(b2)

	var s mysilce = []int{1, 2, 3} // 未命名类型，隐式转换
	var s2 []int = s
	fmt.Println(s2)
}
