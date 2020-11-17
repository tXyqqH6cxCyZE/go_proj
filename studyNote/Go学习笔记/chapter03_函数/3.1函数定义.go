package main

import "fmt"

func test(fn func() int) int {
	return fn()
}

type FormatFunc func(s string, x, y int) string // 定义函数类型 (建议将复杂签名定义为函数类型)

func format(fn FormatFunc, s string, x, y int) string {
	return fn(s, x, y)
}

func main() {

	s1 := test(func() int { return 100 }) // 直接将匿名函数当作参数

	s2 := format(func(s string, x, y int) string {
		return fmt.Sprintf(s, x, y)
	}, "%d,%d", 10, 20)

	println(s1, s2)
}
