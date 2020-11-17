package main

import "fmt"

func test1(s string, n ...int) string {
	var x int
	for _, i := range n {
		x += i
	}
	return fmt.Sprintf(s, x)
}

func main() {
	println(test1("sum: %d", 1, 2, 3))

	// 使用slice对象做变参时，必须展开
	//s := []int{4, 5, 6}
	println(test1("sum: %d", []int{4, 5, 6}...))
}
