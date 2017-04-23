package main

import "net/http"
import "fmt"
import "log"

type mytype struct{}

func (t *mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there from mytype")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func RunSomeCode(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		handler.ServeHTTP(w, r)
		log.Println("After")
		})
}


func main() {

	mux := http.NewServeMux()

	t := new(mytype)
	mux.Handle("/", t)
	mux.HandleFunc("/status/", StatusHandler)

	WrappedMux := RunSomeCode(mux)
	http.ListenAndServe(":8080", WrappedMux)
}
