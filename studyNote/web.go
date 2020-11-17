package studyNote

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()     // analyze the param, dafault : not analyze
	fmt.Print(r.Form) // the info which will output to server
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // this w is which output to client
}

func main() {
	http.HandleFunc("/", sayhelloName)       // set the router of interview
	err := http.ListenAndServe(":9090", nil) // set the port of listening
	if err != nil {
		log.Fatal("ListendAndServer", err)
	}
}
