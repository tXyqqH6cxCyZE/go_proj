package main

import "fmt"

/*
   1. 闭包复制的是原对象指针，这就很容易解释延迟引用现象
   2. 在汇编层面，tmp()实际返回的是FuncVal对象，其中包含匿名函数地址，闭包对象指针。
      当调用匿名函数时，只需以某个寄存器传递该对象即可
      FuncVal { func_address, closure_var_pointer...}
*/
func tmp() func() {
	x := 100
	fmt.Printf("x (%p) = %d\n", &x, x)

	return func() {
		fmt.Printf("x (%p) = %d\n", &x, x)
	}
}

func main() {

	// --- function collection ---
	fns := []func(x int) int{
		func(x int) int { return x + 1 },
		func(x int) int { return x + 2 },
	}
	println(fns[0](100))
	println(fns[1](101))

	// --- function as field ---
	d := struct {
		fn func() string
	}{
		fn: func() string { return "Hello, World!" },
	}
	println(d.fn())

	// --- channel of function ---
	fc := make(chan func() string, 2)
	fc <- func() string { return "Hello, World!" } // 匿名函数通过输入到线道fc
	println((<-fc)())                              // 线道fc输出

	f := tmp()
	f()
}
