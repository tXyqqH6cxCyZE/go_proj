package main

/**
让编译器检查，以确保某个类型实现接口。

某些时候，让函数直接“实现”接口能省下不少事
*/

//type Tester1 interface {
//	DO()
//}
//
//type FuncDo func()
//
//func (self FuncDo) Do() {self()}
//
//func studyNote() {
//	var t Tester1 = FuncDo(func() {
//		println("Hello, World!")
//	})
//	t.DO()
//}
