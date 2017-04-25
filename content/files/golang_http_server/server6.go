// Implementation of the first example in:
// https://github.com/golang/go/issues/18997

package main

import "net/http"
import "fmt"

type statusCaptureWriter struct {
	http.ResponseWriter
	status int
}

func (scw *statusCaptureWriter) WriteHeader(status int) {
	scw.status = status
	scw.ResponseWriter.WriteHeader(status)
}

type loggedHandler struct {
	handler http.Handler
}

func (h *loggedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	scw := &statusCaptureWriter{ResponseWriter: w}
	h.handler.ServeHTTP(scw, r)
	fmt.Println("status", scw.status)
}

type mytype struct{}

func (t *mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// I think I am supposed to call WriteHeader() explicitly?
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there!")
}

func main() {
	t := new(mytype)
	h := &loggedHandler{t}
	http.Handle("/", h)
	http.ListenAndServe(":8080", nil)
}
