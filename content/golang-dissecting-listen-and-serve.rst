:Title: Dissecting golang's HandlerFunc, Handle and DefaultServeMux
:Date: 2017-04-22 10:00
:Category: golang


The `http.ListenAndServe(..) <https://golang.org/pkg/net/http/#ListenAndServe>`__ function is the most straightforward 
approach to start a HTTP 1.1 server. The following code does just that:

.. code-include:: files/golang_http_server/server1.go
    :lexer: python

What is the `nil` second argument above? The documentation states that the second argument to the function should be a "handler" and if it is specified as `nil`, it defaults to `DefaultServeMux`. By the end of next section and eventually by the end of the post, I will have successfully removed one of my main misconceptions - i.e. `DefaultServeMux` was something else other than a handler.


What is `DefaultServeMux`?
==========================

If we run our server above, and send a couple of HTTP GET requests, we will see the following:

.. code::
   
   $ curl localhost:8080
   404 page not found
   
   $ curl localhost:8080/status/
   404 page not found

This is because, we haven't specified how our server should handle requests to GET the root ("/") - our first request or requests to GET the "/status" resource - our second request. Before we see how we could fix that, let's understand *how* the error message "404 page not found" is generated.

The error message is generated from the function below in `src/net/http/server.go` specifically the `NotFoundHandler()` "handler" function:

.. code-include:: files/golang_http_server/snippet1.go
    :lexer: golang


Now, let's roughly see how our GET request above reaches the above function. Let us consider the function signature of the above handler function: `func (mux *ServeMux) handler(host, path string) (h Handler, pattern string)`. 

This function is a method belonging to the type `ServeMux`:

.. code-include:: files/golang_http_server/snippet2.go
    :lexer: golang


So, how does `DefaultServeMux` get set when the second argument to `ListenAndServe()` is `nil`? The following code snippet has the answer:

.. code-include:: files/golang_http_server/snippet3.go
    :lexer: golang


The above call to `ServeHTTP()` calls the following implementation of `ServeHTTP()`:

.. code-include:: files/golang_http_server/snippet4.go
    :lexer: golang

The call to `Handler()` function then calls the following implementation:

.. code-include:: files/golang_http_server/snippet5.go
    :lexer: golang


Now, when we make a request to "/" or "/status/", no match is found by the `mux.match()` call above and hence the handler returned is the `NotFoundHandler` whose `ServeHTTP()` function is then called to return the "404 page not found" error message:

.. code-include:: files/golang_http_server/snippet6.go
    :lexer: golang

We will discuss how this "magic" happens in the next section.

Registering handlers
====================

Let's now update our server code to handle "/" and "/status/":

.. code-include:: files/golang_http_server/server2.go
    :lexer: golang

If we run the server and send the two requests above, we will see the following responses:

.. code::

   $ curl localhost:8080
   Hello there from mytype 

   $ curl localhost:8080/status/
   OK


Let's look at what the call to `Handle()` function does:

.. code-include:: files/golang_http_server/snippet7.go
    :lexer: golang

Now, let's look at what the call to `HandleFunc()` function does:

.. code-include:: files/golang_http_server/snippet8.go
    :lexer: golang

The call to the ``http.HandleFunc()`` function "converts" the provided function to the ``HandleFunc()`` type and then calls the ``(mux *ServeMux) Handle()`` function similar to what happens when we call the ``Handle()`` function.

Let's now revisit how the right handler function gets called. In a code snippet above, we saw a call to the ``match()`` function which given a path returns the most appropriate registered handler for the path:


.. code-include:: files/golang_http_server/snippet9.go
    :lexer: golang

``mux.m`` is a a ``map`` data structure defined in the ``ServeMux`` structure (snippet earlier in the post).

The HandleFunc() type
=====================


Using your own Handler with ListenAndServe()
============================================


Writing Middleware
==================


Learn more
==========

- http://jordanorelli.com/post/42369331748/function-types-in-go-golang
- https://golang.org/doc/effective_go.html#interface_methods












