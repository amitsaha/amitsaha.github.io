:Title: Dissecting golang's HandlerFunc, Handle and DefaultServeMux
:Date: 2016-04-17 10:00
:Category: golang 

The `http.ListenAndServe(..) <https://golang.org/pkg/net/http/#ListenAndServe>`__ function is the most straightforward 
approach to start a HTTP 1.1 server. The following code does just that:

.. code-include: files/golang_http_server/server1.py
    :lexer: python

**Question 1**: What is the `nil` second argument above? If not specified, it defaults to `DefaultServeMux`. What is it? We will find out soon.

Now, if we try to send a couple of HTTP GET requests, we will see the following:

.. code::
   
   $ curl localhost:8080
   404 page not found
   
   $ curl localhost:8080/status/
   404 page not found

This is because, we haven't specified how our server should handle requests to GET the root ("/") - our first request or requests to GET the "/status" resource - our second request. Before we see how we could fix that, let's understand *how* the error message "404 page not found" is generated.

The error message is generated from the function below in `src/net/http/server.go` specifically the `NotFoundHandler()` "handler" function:

.. code-inclide: files/golang_http_server/snippet1.go
    :lexer: golang


We will discuss handlers shortly. First, let's consider the function signature of the above handler function: `func (mux *ServeMux) handler(host, path string) (h Handler, pattern string)`. This function is a method belonging to the type `ServeMux`:

.. code-inclide: files/golang_http_server/snippet2.go
    :lexer: golang


**Answer 1**:  Let's now go back to our question 1 we asked: How does `DefaultServeMux` get set when the second argument is `nil`? The following code snippet has the answer:

.. code-inclide: files/golang_http_server/snippet3.go
    :lexer: golang


The above call to `ServeHTTP()` calls the following implementation of `ServeHTTP()`:

.. code-inclide: files/golang_http_server/snippet4.go
    :lexer: golang

The call to `Handler()` function then calls the following implementation of `Handler()`:

.. code-inclide: files/golang_http_server/snippet5.go
    :lexer: golang


Now, when we make a request to "/" or "/status/", no match is found by the `mux.match()` call above and hence the handler returned is the `NotFoundHandler` whose `ServeHTTP()` function is then called to return the "404 page not found" error message:

.. code-inclide: files/golang_http_server/snippet6.go
    :lexer: golang









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




