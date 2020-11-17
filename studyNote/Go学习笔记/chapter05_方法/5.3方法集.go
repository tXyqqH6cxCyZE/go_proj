package main

/**
方法集 ？
	每个类型都有与之关联的方法集，这会影响到接口实现规则。
		- 类型T方法集包含全部 receiver T 方法。
		- 类型*T方法集合包含全部 receiver T + *T 方法。
		- 如类型S包含匿名字段T，则S方法集包含T方法。
		- 如类型S包含匿名字段*T,则S方法集包含 T + *T 方法。
		- 不管嵌入T 或 *T,*S方法集总是包含T + *T方法。
用实例value和pointer调用方法（含匿名字段）不受方法集约束，
编译器总是查找全部方法，并自动转换receiver实参
*/
