package main

import "net/http"
import "fmt"
import "log"


type MyResponseWriter struct {
	http.ResponseWriter
	code       int
}

/*func (mw *MyResponseWriter) Write(p []byte) (int, error) {
	mw.WriteHeader(http.StatusOK)
	return mw.ResponseWriter.Write(p)
}
*/

func (mw *MyResponseWriter) WriteHeader(code int) {
	mw.code = code
	mw.ResponseWriter.WriteHeader(code)
}

type mytype struct{}

func (t *mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeHTTP called")
	fmt.Fprintf(w, "Hello there from mytype")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func RunSomeCode(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a %s request for: %v", r.Method, r.URL)
		myrw := &MyResponseWriter{ResponseWriter: w}
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
