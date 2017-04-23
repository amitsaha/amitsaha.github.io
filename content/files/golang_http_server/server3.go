package main

import "net/http"
import "fmt"

type mytype struct{}

func (t *mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there from mytype")
}


func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {

	t := new(mytype)
	http.Handle("/", t)
	
	http.HandleFunc("/status/", StatusHandler)
        
	http.ListenAndServe(":8080", nil)
}
