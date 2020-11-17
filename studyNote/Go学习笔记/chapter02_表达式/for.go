package main

import "fmt"

func length(s string) int {
	println("call length.")
	return len(s)
}
func main() {

	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		println(k, v)
	}

	/**
	注意：range会复制对象
	*/
	a := [3]int{0, 1, 2}

	for i, v := range a { // index, value 都是从复制品中取出
		if i == 0 { // 在修改前，我们先修改原数组
			a[1], a[2] = 999, 99
			fmt.Println(a) // 确认修改有效，输出[0,999,999]
		}

		a[i] = v + 100 // 使用复制品中取出的value修改原数组
	}
	fmt.Println(a) // 输出[100, 101, 102].
	fmt.Println("-------------------------------------------------------")
	/**
	建议使用引用类型，其底层数据不会被复制
	注意：另外两种引用类型map，channel是指针包装，而不像slice是struct
	*/
	s := []int{1, 2, 3, 4, 5}
	for i, v := range s { // 复制 struct slice {pointer, len, cap}
		if i == 0 {
			s = s[:3]  // 对slice的修改，不会影响range
			s[2] = 100 // 对底层数据的修改
		}

		println(i, v)
	}

}
