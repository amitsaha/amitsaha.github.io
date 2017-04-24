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
        log.Printf("Got a %s request for: %v", r.Method, r.URL)
		handler.ServeHTTP(w, r)
        // At this stage, our handler has "handled" the request
        // but we can still write to the client there
        // but we won't do that
        // XXX: We don't have the HTTP status here either, need to understand this better why
        log.Println("Handler finished processing request")
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
