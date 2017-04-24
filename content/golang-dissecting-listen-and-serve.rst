:Title: Dissecting golang's HandlerFunc, Handle and DefaultServeMux
:Date: 2017-04-24 10:00
:Category: golang

My aim in this post is to discuss three "concepts" in Golang that I came across while writing HTTP servers. Through this
discussion I am to get rid of my own lack of understanding (at least to a certain degree) about these. Hopefully, it will
be of use to others too.

The `http.ListenAndServe(..) <https://golang.org/pkg/net/http/#ListenAndServe>`__ function is the most straightforward 
approach to start a HTTP 1.1 server. The following code does just that:

.. code-include:: files/golang_http_server/server1.go
    :lexer: python

What is the `nil` second argument above? The documentation states that the second argument to the function should be a 
"handler" and if it is specified as `nil`, it defaults to `DefaultServeMux`.


What is `DefaultServeMux`?
==========================

If we run our server above via ``go run server1.go``, and send a couple of HTTP GET requests, we will see the following:

.. code::
   
   $ curl localhost:8080
   404 page not found
   
   $ curl localhost:8080/status/
   404 page not found

This is because, we haven't specified how our server should handle requests to GET the root ("/") - our first request or 
requests to GET the "/status" resource - our second request. Before we see how we could fix that, let's understand 
*how* the error message "404 page not found" is generated.

The error message is generated from the function below in `src/net/http/server.go` specifically the `NotFoundHandler()` 
"handler" function:

.. code-include:: files/golang_http_server/snippet1.go
    :lexer: go


Now, let's roughly see how our GET request above reaches the above function. 

Let us consider the function signature of the above handler function: `func (mux *ServeMux) handler(host, path string) (h Handler, pattern string)`. This function is a method belonging to the type `ServeMux`:

.. code-include:: files/golang_http_server/snippet2.go
    :lexer: go


So, how does `DefaultServeMux` get set when the second argument to `ListenAndServe()` is `nil`? The following code 
snippet has the answer:

.. code-include:: files/golang_http_server/snippet3.go
    :lexer: go


The above call to `ServeHTTP()` calls the following implementation of `ServeHTTP()`:

.. code-include:: files/golang_http_server/snippet4.go
    :lexer: go

The call to `Handler()` function then calls the following implementation:

.. code-include:: files/golang_http_server/snippet5.go
    :lexer: go


Now, when we make a request to "/" or "/status/", no match is found by the `mux.match()` call above and hence the 
handler returned is the `NotFoundHandler` whose `ServeHTTP()` function is then called to return the "404 page not found" 
error message:

.. code-include:: files/golang_http_server/snippet6.go
    :lexer: go

We will discuss how this "magic" happens in the next section.

Registering handlers
====================

Let's now update our server code to handle "/" and "/status/":

.. code-include:: files/golang_http_server/server2.go
    :lexer: go

If we run the server and send the two requests above, we will see the following responses:

.. code::

   $ curl localhost:8080
   Hello there from mytype 

   $ curl localhost:8080/status/
   OK



Let's now revisit how the right handler function gets called. In a code snippet above, we saw a call to the ``match()`` function which given a path returns the most appropriate registered handler for the path:


.. code-include:: files/golang_http_server/snippet9.go
    :lexer: go

``mux.m`` is a a ``map`` data structure defined in the ``ServeMux`` structure (snippet earlier in the post) which stores a mapping of a path and the handler we have registered for it.

**The HandleFunc() type**

Let's go back to the idea of "converting" any function with the signature ``func aFunction(w http.ResponseWriter, r *http.Request)`` to the type "HandlerFunc". 

Any type which has a ServeHTTP() method is said to implement the ``Handler`` interface:

.. code::

    type HandlerFunc func(ResponseWriter, *Request)

    // ServeHTTP calls f(w, req).
    func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
        f(w, req)
    }


Going back to the previous version of our server, we see how we do that:


.. code::

    type mytype struct{}

    func (t *mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello there from mytype")
    }

The ``ServeHTTP()`` method of a Handler is invoked when it has been registered as handling a particular path.

Let's look at what the call to `Handle()` function does:

.. code-include:: files/golang_http_server/snippet7.go
    :lexer: go


It can feel cumbersome to define a type implementing the ``Handler`` interface for every path we want to register a handler for. Hence, a convenience function, ``HandleFunc()`` is provided to register any function which has a specified signature as a Handler function. For example:

.. code::

    http.HandleFunc("/status/", StatusHandler)

Now, let's look at what the call to `HandleFunc()` function does:

.. code-include:: files/golang_http_server/snippet8.go
    :lexer: go

The call to the ``http.HandleFunc()`` function "converts" the provided function to the ``HandleFunc()`` type and then calls the ``(mux *ServeMux) Handle()`` function similar to what happens when we call the ``Handle()`` function. The idea of this conversion is explained in the `Effective Go guide <https://golang.org/doc/effective_go.html#interface_methods>`__ and this `blog post <http://jordanorelli.com/post/42369331748/function-types-in-go-golang>`__.



Using your own Handler with ListenAndServe()
============================================

Earlier in this post, we saw how passsing ``nil`` to ``ListenAndServe()`` function sets the handler to ``DefaultServeMux``. The handlers
we register via ``Handle()`` and ``HandleFunc()`` are then added to this object. Hence, we could without changing any functionality rewrite our server as follows:

.. code-include:: files/golang_http_server/server3.go
    :lexer: go

We create an object of type ``ServeMux`` via ``mux := http.NewServeMux()``, register our handlers calling the same two functions, but those that are defined for the ``ServeMux`` object we created.

The reason we may want to use our own Handler with ``ListenAndServe()`` is demonstrated in the next section.


Writing Middleware
==================

In our latest version of the server, we have specified our own handler to ``ListenAndServe()``. One reason for doing so is when you want to execute some code for every request. That is:

1. Server gets a request for "/path/"
2. Execute some code
3. Handler for "/path/" gets called
4. Execute some code
5. Return the response to the client

Either of steps 2 or 4 or both may occur and this is where "middleware" comes in. Our next version of the server demonstrates how we may implement this:


.. code-include:: files/golang_http_server/server4.go
    :lexer: go

When we run the server and send it a couple of requests as above, we will see:

.. code::

    2017/04/24 17:53:03 Got a GET request for: /
    2017/04/24 17:53:03 Handler finished processing request
    2017/04/24 17:53:05 Got a GET request for: /status
    2017/04/24 17:53:05 Handler finished processing request

As part of this middleware writing exercise, I also wanted to be able to print the HTTP status of the response that we are sending but as the comment in the code states,
I haven't been able to figure it out yet.


Learn more
==========

- https://gocodecloud.com/blog/2016/11/15/simple-golang-http-request-context-example/
- http://jordanorelli.com/post/42369331748/function-types-in-go-golang
- https://golang.org/doc/effective_go.html#interface_methods












