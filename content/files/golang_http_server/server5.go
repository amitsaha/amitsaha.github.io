package main

import "net/http"
import "fmt"
import "log"

type MyResponseWriter struct {
	http.ResponseWriter
	code int
}

func (mw *MyResponseWriter) Header() http.Header {
	return mw.ResponseWriter.Header()
}

func (mw *MyResponseWriter) WriteHeader(code int) {
	mw.code = code
	mw.ResponseWriter.WriteHeader(code)
}

type mytype struct{}

func (t *mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there from mytype")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func RunSomeCode(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a %s request for: %v", r.Method, r.URL)
		myrw := &MyResponseWriter{ResponseWriter: w, code: -1}
		handler.ServeHTTP(myrw, r)
		log.Println("Response status: ", myrw.code)
	})
}

func main() {

	mux := http.NewServeMux()

	t := new(mytype)
	mux.Handle("/", t)
	mux.HandleFunc("/status/", StatusHandler)

	WrappedMux := RunSomeCode(mux)
	log.Fatal(http.ListenAndServe(":8080", WrappedMux))
}
