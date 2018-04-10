package main

import (
"fmt"
"log"
"net/http"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["gm"])
	/*for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}*/
	//fmt.Fprintf(w, "Hello world!")
}
func main() {
	http.HandleFunc("/", sayhelloName)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
//for example
//http://localhost:8080/url_long?12344=1&123456=1222
