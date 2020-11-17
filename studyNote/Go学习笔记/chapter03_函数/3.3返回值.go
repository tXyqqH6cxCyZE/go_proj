package main

func test2() (int, int) {
	return 1, 2
}

func add(x, y int) int {
	return x + y
}

func sum(n ...int) int {
	var x int
	for _, i := range n {
		x += i
	}
	return x
}

func add1(x, y int) (z int) {
	z = x + y
	return
}

func add2(x, y int) (z int) {
	defer func() {
		z += 100
	}()

	z = x + y
	return
}

func add3(x, y int) (z int) {
	defer func() {
		println(z)
	}()

	z = x + y
	return z + 200
}

func main() {
	// 多返回值可直接作为其他函数调用实参
	println(add(test2()))
	println(sum(test2()))
	/**
	  1. 命名返回参数可看做与形参类似的局部变量，最后由return隐式返回
	  2. 命名返回参数会被同名局部变量遮蔽，此时需要显式返回
	  3. 命名返回参数允许defer延迟调用通过闭包读取和修改
	  4. 显式return返回前，会先修改命名返回参数
	*/
	println(add1(1, 2))

	println(add2(1, 2))

	println(add3(1, 2))
}
