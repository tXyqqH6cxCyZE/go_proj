package main

import "fmt"

/**
需要说明，slice并不是数组或数组指针。它通过内部指针和相关属性引用数组片段，以实现变长方案
runtime.h
-------------------------------------------
struct Slice
{							// must not move anything
	byte* array;			// actual data
	uintgo	len;			// number of elements
	uintgo	cap;			// allocated number of elements
}

- 引用类型。但自身是结构体，值拷贝传递。
- 属性 len 表示可用元素数量，读写操作不能超过该限制。
- 属性 cap 表示最大扩张容量，不能超出数组限制。
- 如果 slice == nil, 那么len, cap都等于0
*/

func main() {

	data := []int{0, 1, 2, 3, 4, 5, 6}
	slice := data[1:4:5]

	fmt.Println(slice)
	fmt.Println(cap(slice))
	fmt.Println(slice[2])
	fmt.Println("-------------------------------------")

	/**
	读写操作实际目标是底层数组，只需注意索引号的差别。
	*/
	dat := []int{0, 1, 2, 3, 4, 5}
	s := data[2:4]
	s[0] += 100
	s[1] += 200

	fmt.Println(s)
	fmt.Println(dat)
	fmt.Println("-------------------------------------")

	/**
	使用 make 动态创建slice,避免了数组必须用常量做长度的麻烦。
	还可以用指针直接访问底层数组，退化成普通数组操作
	*/
	ts := []int{0, 1, 2, 3}
	p := &ts[2] // *int, 获取底层数组元素指针
	*p += 100
	fmt.Println(ts)

	/**
		1 reslice:
	 	- 所谓reslice，是基于已有slice创建新slice对象，以便在cap允许范围内调整属性。
		- 新对象依然指向原底层数组
	*/
	rs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := rs[2:5] // [2, 3, 4]
	s1[2] = 100

	s2 := s1[2:6] // [100 5 6 7]
	s2[3] = 200
	fmt.Println(rs)

	/**
	2 append
		2.1 向slice尾部添加数据，返回新的slice对象
		2.2 简单地说，就是在array[slice.high]写数据
		2.3 一旦超出原slice.cap限制，就会重新分配底层数组，即使原数组并未填满
		2.4 通常以2倍容量重新分配数组。在大批量添加数据时，建议一次性分配足够大的空间，以减少
			内存分配和数据复制开销。或初始化足够长的len属性，改用索引号进行操作。
			及时释放不再使用的slice对象，避免持有过期数组，造成GC无法回收。
	*/
	// 2.1
	as := make([]int, 0, 5)
	fmt.Printf("%p\n", &as)

	as2 := append(as, 1)
	fmt.Printf("%p\n", &as2)

	fmt.Println(as, as2)

	// 2.3
	adata := [...]int{0, 1, 2, 3, 4, 10: 0}
	as3 := adata[:2:3]

	as3 = append(as3, 100, 200) // 一次append两个值，超出as3.cap限制

	fmt.Println(as3, adata)         // 重新分配底层数组，与原数组无关
	fmt.Println(&as3[0], &adata[0]) // 比对底层数组起始指针，发现内存地址不同

	/**
		3 copy
		- 函数copy在两个slice间复制数据，复制长度以len小的为准。
	      两个slice可指向同一底层数组，允许元素区间重叠。
		- 应及时将所需数据copy到较小的slice，以便释放超大号底层数组内存
	*/
	cdata := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	cs := cdata[8:]
	cs2 := cdata[:5]

	copy(cs2, cs)

	fmt.Println(cs2)
	fmt.Println(cdata)
}
