package main

import (
	"db/dao"
	json2 "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/**
获取请求体中的信息
*/

// 普通响应
func handler(w http.ResponseWriter, r *http.Request) {
	// 获取内容的长度
	length := r.ContentLength
	// 创建一个字节切片
	body := make([]byte, length)
	// 读取请求体
	r.Body.Read(body)
	fmt.Fprintln(w, "请求体中的内容是：", string(body))
}

// 响应文件
func handlerGetFile(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	r.ParseMultipartForm(1024)

	fileHeader := r.MultipartForm.File["text"][0]

	file, err := fileHeader.Open()

	// net/http提供的FormFile方法可以快速的获取被上传的文件
	// 但是只能处理上传一个文件的情况
	//file, _, err := r.FormFile("photo")

	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

// 响应json数据
func handlerJson(w http.ResponseWriter, r *http.Request) {
	// 设置响应头中内容的类型
	w.Header().Set("Content-Type", "application/json")
	user := dao.User{
		ID:       1,
		Username: "admin",
		Password: "123456",
		Email:    "pjx@xiaomi.com",
	}

	// 将user转换为json格式
	json, _ := json2.Marshal(user)
	w.Write(json)
}

// 让客户端重定向
func handlerDirect(w http.ResponseWriter, r *http.Request) {
	// 以下操作必须要在WriteHeader之前进行
	w.Header().Set("Location", "https:www.baidu.com")
	w.WriteHeader(302)
}

// cookie机制
func handlerCookie(w http.ResponseWriter, r *http.Request) {
	cookie1 := http.Cookie{
		Name:     "user1",
		Value:    "admin",
		HttpOnly: true,
		MaxAge:   60, // 设置Cookie有效时间
	}

	cookie2 := http.Cookie{
		Name:     "user2",
		Value:    "superAdmin",
		HttpOnly: true,
		MaxAge:   60, // 设置Cookie有效时间
	}

	// 将Cookie发送给浏览器，即添加第一个Cookie
	w.Header().Set("Set-Cookie", cookie1.String())
	// 再添加一个Cookie
	w.Header().Add("Set-Cookie", cookie2.String())

	// 除了Set和Add方法之外，Go还提供了一种更快捷的设置Cookie的方式，
	// 就是通过net/http库中的SetCookie方法
	//http.SetCookie(w, &cookie1)
	//http.SetCookie(w, &cookie2)
}

func handlerGetCookies(w http.ResponseWriter, r *http.Request) {
	// 获取请求头中的Cookie
	cookies := r.Header["Cookie"]
	fmt.Fprintln(w, cookies)
}

func main() {
	http.HandleFunc("/getBody", handler)
	http.HandleFunc("/upload", handlerGetFile)
	http.HandleFunc("/json", handlerJson)
	http.HandleFunc("/direct", handlerDirect)
	http.HandleFunc("/cookie", handlerCookie)
	http.HandleFunc("/getCookies", handlerGetCookies)

	// ListendAndServe要放在最后
	http.ListenAndServe(":8080", nil)
}
