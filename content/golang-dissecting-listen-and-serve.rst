:Title: Dissecting golang's HandlerFunc, Handle and DefaultServeMux
:Date: 2017-04-22 10:00
:Category: golang 

The `http.ListenAndServe(..) <https://golang.org/pkg/net/http/#ListenAndServe>`__ function is the most straightforward 
approach to start a HTTP 1.1 server. The following code does just that:

.. code-include:: files/golang_http_server/server1.go
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

.. code-include:: files/golang_http_server/snippet1.go
    :lexer: golang


We will discuss handlers shortly. First, let's consider the function signature of the above handler function: `func (mux *ServeMux) handler(host, path string) (h Handler, pattern string)`. This function is a method belonging to the type `ServeMux`:

.. code-include:: files/golang_http_server/snippet2.go
    :lexer: golang


**Answer 1**:  Let's now go back to our question 1 we asked: How does `DefaultServeMux` get set when the second argument is `nil`? The following code snippet has the answer:

.. code-include:: files/golang_http_server/snippet3.go
    :lexer: golang


The above call to `ServeHTTP()` calls the following implementation of `ServeHTTP()`:

.. code-include:: files/golang_http_server/snippet4.go
    :lexer: golang

The call to `Handler()` function then calls the following implementation of `Handler()`:

.. code-include:: files/golang_http_server/snippet5.go
    :lexer: golang


Now, when we make a request to "/" or "/status/", no match is found by the `mux.match()` call above and hence the handler returned is the `NotFoundHandler` whose `ServeHTTP()` function is then called to return the "404 page not found" error message:

.. code-include:: files/golang_http_server/snippet6.go
    :lexer: golang











