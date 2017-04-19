The `http.ListenAndServe(..) <https://golang.org/pkg/net/http/#ListenAndServe>`__ function is the 




// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, req).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
    f(w, req)
}


// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.Handle(pattern, HandlerFunc(handler))
}

// Handle registers the handler for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}



// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern " + pattern)
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if mux.m[pattern].explicit {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	mux.m[pattern] = muxEntry{explicit: true, h: handler, pattern: pattern}
...

// Find a handler on a handler map given a path string
// Most-specific (longest) pattern wins
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	var n = 0
	for k, v := range mux.m {
		if !pathMatch(k, path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = v.h
			pattern = v.pattern
		}
	}
	return
}





// Handler returns the handler to use for the given request,
// consulting r.Method, r.Host, and r.URL.Path. It always returns
// a non-nil handler. If the path is not in its canonical form, the
// handler will be an internally-generated handler that redirects
// to the canonical path.
//
// Handler also returns the registered pattern that matches the
// request or, in the case of internally-generated redirects,
// the pattern that will match after following the redirect.
//
// If there is no registered handler that applies to the request,
// Handler returns a ``page not found'' handler and an empty pattern.
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path {
			_, pattern = mux.handler(r.Host, p)
			url := *r.URL
			url.Path = p
			return RedirectHandler(url.String(), StatusMovedPermanently), pattern
		}
	}

	return mux.handler(r.Host, r.URL.Path)
}

// handler is the main implementation of Handler.
// The path is known to be in canonical form, except for CONNECT methods.
func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path)
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}


// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.logf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		if !c.hijacked() {
			c.close()
			c.setState(c.rwc, StateClosed)
		}
	}()

....

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(rw, req)
}


ListenAndServe(): nil 

DefaultServeMux is actually another Handler

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type Server struct {
	Addr      string      // TCP address to listen on, ":http" if empty
	Handler   Handler     // handler to invoke, http.DefaultServeMux if nil
	TLSConfig *tls.Config // optional TLS config, used by ListenAndServeTLS
...

ListenAndServe starts an HTTP server with a given address and handler. The handler is usually nil, which means to use DefaultServeMux. Handle and HandleFunc add handlers to DefaultServeMux:


Example code:

package main

import "net/http"
import "fmt"

type mytype struct{}

func (t *mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there from Handlefunc")
}

func main() {

	t := new(mytype)
	http.Handle("/", t)
	http.ListenAndServe(":8080", t)
}


Anything that has a ServeHTTP() method can be used as a handler.

// HTTP server with middleware
package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Send a response
	fmt.Fprintf(w, "Hi there! %s", r.URL.Path[1:])
	// Send another response
	fmt.Fprintf(w, "Hi again!")
}

// This handles /ping requests, but is wrapped by the "makeHandler"
// function below
func pingHandler(w http.ResponseWriter, r *http.Request, s string) {
	fmt.Fprintf(w, "Pong! %s", s)
}

// This is a function which returns a http.HandlerFunc, useful to implement common logic
// across Handlers
// It accepts a function as argument which takes in http.ResponseWriter, *http.Request
// and a string
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, "testing")
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ping", makeHandler(pingHandler))
	// nil argument here specifies using the DefaultServeMux
	http.ListenAndServe(":8080", makeHandler(pingHandler))
}




HandleFunc()
============

func rootHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(w, "Hello there from Handlefunc")
}

http.HandleFunc("/", rootHandler)




Handle()
========

type mytype struct {}

func (t mytype) ServeHTTP(w ResponseWriter, req *Request) {
	fmt.Printf(w, "Hello there from Handle")
}

t = new(mytype)
http.Handle("/custom", t)




http://jordanorelli.com/post/42369331748/function-types-in-go-golang
https://golang.org/doc/effective_go.html#interface_methods




