package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// if we not use the ParseForm(), we will not get the data below
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello dwn!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		// if the data is login data, then judge the logic of login
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		//fmt.Println(r.Form)
	}
}

func main() {
	http.HandleFunc("/", sayhelloName)       // set the router of interview
	http.HandleFunc("/login", login)         // set the router of interview
	err := http.ListenAndServe(":9090", nil) // set the port of listening
	if err != nil {
		log.Fatal("ListenAndServer", err)
	}
}
