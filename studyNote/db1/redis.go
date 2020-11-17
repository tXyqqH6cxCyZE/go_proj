package main

import (
	"fmt"
	"github.com/astaxie/goredis"
)

func main() {
	var client goredis.Client
	// the operation of string
	client.Set("a", []byte("hello"))
	// 为什么
	val, _ := client.Get("a")
	fmt.Println(string(val))
	//client.Del("a")

	// the operation of list
	vals := []string{"a", "b", "c", "d", "e"}
	for _, v := range vals {
		client.Rpush("1", []byte(v))
	}

	dbvals, _ := client.Lrange("l", 0, 4)
	for i, v := range dbvals {
		println(i, ":", string(v))
	}

	client.Del("l")
}
