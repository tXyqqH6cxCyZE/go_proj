package main

import "fmt"

/**
没有结构化异常，使用panic抛出错误，recover捕获错误
*/

/**
由于painc,recover 参数类型 为interface{}, 因此可抛出任何类型对象
func panic(v interface{})
func recover() interface{}
*/
func test9() {
	defer func() {
		if err := recover(); err != nil {
			println(err.(string)) // 将 interface{} 转型为具体类型
		}
	}()
	panic("panic error!")
}

/**
延迟调用中引发的错误，可被后续延迟调用捕获，但仅最后一个错误可被捕获
*/
func test10() {
	defer func() {
		fmt.Println(recover())
	}()

	defer func() {
		panic("defer panic")
	}()

	panic("test panic")
}

/**
捕获函数 recover 只有在延迟调用内直接调用才会终止错误，否则总是返回nil。
任何未捕获的错误都会沿调用堆栈向外传递。
*/
func test11() {
	defer recover()              // 无效！
	defer fmt.Println(recover()) // 无效！
	defer func() {
		func() {
			println("defer inner")
			recover() // 无效！
		}()
	}()

	panic("test panic")
}

func main() {
	//test9()
	//test10()
	test11()
}
