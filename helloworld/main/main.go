// Autogenerated by Thrift Compiler (0.9.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"/home/ubuntu/proj/helloworld/framework/app"
	"/home/ubuntu/proj/helloworld/hello"
	"fmt"
)

func main() {
	handler := NewHelloServiceHandler()
	processor := hello.NewHelloServiceProcessor(handler)
	hook := &XMainHook{processor}
	err := app.Run(processor, hook)
	if err != nil {
		fmt.Println("run app failed, err:", err)
	}
	hook.OnShutdown()
}
