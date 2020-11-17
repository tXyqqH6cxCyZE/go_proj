package main

/**
	Map:
	- 引用类型，哈系表。键必须是支持相等运算符（==，!=）类型，比如number，string，
  	pointer,array,struct,以及对应的interface。值可以是任意类型，没有限制。
	- 可以在迭代时安全删除键值。但如果期间有新增操作，那么就不知道会有什么意外了
*/

func main() {

	// 从map中取回的是一个value复制品，对其成员的修改是没有意义的

	type user struct{ name string }

	m := map[int]user{ // 当map因扩张而重新哈希时，各键值项存储位置都会发生改变。因此，map
		1: {"user1"}, // 被设计成 not addressable。类似m[1].name这种期望透过原value
	} // 指针修改成员的行为自然会被禁止

	//m[1].name = "Tom"			Error: cannot assign to m[1].name

	// 正确的做法是完整替换value或使用指针
	u := m[1]
	u.name = "Tom"
	m[1] = u // 替换value

	m2 := map[int]*user{
		1: &user{"user1"},
	}
	m2[1].name = "Jack" // 返回的是指针复制品。透过指针修改原对象是允许的
}
