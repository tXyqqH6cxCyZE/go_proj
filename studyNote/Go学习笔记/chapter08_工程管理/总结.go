package main

import (
	"fmt"
	"time"
)

/**
8.3 包结构
	- 源文件头部以"package<name>"声明包名称。
	- 包由同一目录下的多个源码文件组成。
	- 包名类似namespace,与包所在目录名，编译文件名无关。
	- 可执行文件必须包含 package studyNote,入口函数 studyNote。
*/

/**
8.3.2 初始化
	- 每个源文件都可以定义一个或多个初始化函数。
	- 编译器不保证多个初始化函数执行次序。
	- 初始化函数在单一线程被调用，仅执行一次。
	- 初始化函数在包所有全局变量初始化后执行。
	- 在所有初始化函数结束后才执行main.studyNote
	- 无法调用初始化函数

1. 因为无法保证初始化函数执行顺序，因此全局变量应该直接用var初始化
2. 可在初始化函数中使用goroutine，可等待其结束。
3. 不应该滥用初始化函数，仅适合完成当前文件中的相关环境设置。
*/

var now = time.Now()

func init() {
	fmt.Printf("now: %v\n", now)
}

func init() {
	fmt.Printf("since: %v\n", time.Now().Sub(now))
}
