package main

import (
	"os"
	"sync"
	"testing"
)

func test5() error {
	f, err := os.Create("test.txt")
	if err != nil {
		return err
	}

	defer f.Close() // 注册调用，而不是注册函数。必须提供参数，哪怕为空

	f.WriteString("Hello, World!")
	return nil
}

/**
多个defer注册，按FILO次序执行。哪怕函数或某个延迟调用发生错误，这些调用依旧会被执行。
*/
func test6(x int) {
	defer println("a")
	defer println("b")

	defer func() {
		println(100 / x)
	}()

	defer println("c")
}

/**
延迟调用参数在注册时求值或复制，可用指针或闭包“延迟”读取。
*/
func test7() {
	x, y := 10, 20

	defer func(i int) {
		println("defer:", i, y) // y 闭包引用
	}(x) // x 被复制

	x += 10
	y += 100
	println("x =", x, "y =", y)
}

/**
滥用defer可能会导致性能问题，尤其是在一个“大循环”里
*/
var lock sync.Mutex

func test8() {
	lock.Lock()
	lock.Unlock()
}

func testdefer() {
	lock.Lock()
	defer lock.Unlock()
}

func BenchmarkTest(b *testing.B) { // 单元测试
	for i := 0; i < b.N; i++ {
		test8()
	}
}

func BenchmarkTestDefer(b *testing.B) { // 单元测试
	for i := 0; i < b.N; i++ {
		testdefer()
	}
}

func main() {
	//test6(0)
	//test7()
}
