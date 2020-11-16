package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)
/**
 	获取请求体中的信息
 */

func handler(w http.ResponseWriter, r *http.Request){
	// 获取内容的长度
	length := r.ContentLength
	// 创建一个字节切片
	body := make([]byte, length)
	// 读取请求体
	r.Body.Read(body)
	fmt.Fprintln(w, "请求体中的内容是：", string(body))
}

func handlerGetFile(w http.ResponseWriter, r *http.Request){
	// 解析表单
	r.ParseMultipartForm(1024)

	fileHeader := r.MultipartForm.File["text"][0]

	file, err := fileHeader.Open()

	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

func main(){
	http.HandleFunc("/getBody", handler)
	http.HandleFunc("/upload", handlerGetFile)


	// ListendAndServe要放在最后
	http.ListenAndServe(":8080", nil)
}
