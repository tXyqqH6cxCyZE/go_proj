package main

import "fmt"

type Vertex struct {
	X, Y int
}

var (
	v1 = Vertex{1, 2} // 创建一个Vertex类型的结构体
	v2 = Vertex{X: 1} // Y:0 被隐式的赋予
	v3 = Vertex{}     // X:0 Y:0
	p  = &Vertex{1, 2}
)

func main() {
	fmt.Println(v1, p, v2, v3)
}
