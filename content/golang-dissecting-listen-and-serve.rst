:Title: Dissecting golang's HandlerFunc, Handle and DefaultServeMux
:Date: 2017-04-24 10:00
:Category: golang

My aim in this post is to discuss three "concepts" in Golang that I come across while writing HTTP servers. Through this
post, my aim to get rid of my own lack of understanding (at least to a certain degree) about these. Hopefully, it will
be of use to others too. The code references are from `src/net/http/server.go <https://golang.org/src/net/http/server.go>`__. 

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
requests to GET the "/status/" resource - our second request. Before we see how we could fix that, let's understand 
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

In our latest version of the server, we have specified our own handler to ``ListenAndServe()``. One reason for doing so is when you want to execute some code for *every* request. That is:

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

What we are doing above is we are "wrapping" our actual handler in another function ``RunSomeCode(handler http.Handler) http.Handler`` which satisfies the ``Handler`` interface. In this function, we print a log message, then call the ``ServeHTTP()`` method of our original
handler, ``mux``. Once it returns from there, we are then printing another log message.

As part of this middleware writing exercise, I also wanted to be able to print the HTTP status of the response that we are sending but as the comment in the code states, there is no direct way to get the status via the ``ResponseWriter`` object. Our next server example will fix this.

Rewrapping ``http.ResponseWriter``
==================================

It took me a while to write the next version of the server, and after reading through some mailing list postings and example code, 
i have a version which achieves what I wanted to be able to do via my middleware:

.. code-include:: files/golang_http_server/server5.go
    :lexer: go


In the example above, I define a new type ``MyResponseWriter`` which implements the ``http.ResponseWriter`` interface by implementing the
three methods ``Header()``, ``Write()`` and ``WriteHeader()``. In bothe ``Write()`` and ``WriteHeader()``, I have some custom code that I execute before calling the corresponding method defined on the ``http.ResponseWriter()`` interface. 


Then, in ``RunSomeCode()``, instead of using the standard ``http.ResponseWriter()`` object that it was passed, I wrap it in a ``MyResponseWriter`` type as follows:

.. code::
    
    myrw := &MyResponseWriter{ResponseWriter: w, code: -1}
    handler.ServeHTTP(myrw, r)


Now, if we run the server, we will see log messages on the server as follows when we send it HTTP get requests:

.. code::

    2017/04/25 17:33:06 Got a GET request for: /status/
    2017/04/25 17:33:06 Response status:  200
    2017/04/25 17:33:07 Got a GET request for: /status
    2017/04/25 17:33:07 Response status:  301
    2017/04/25 17:33:10 Got a GET request for: /
    2017/04/25 17:33:10 Response status:  200


I will end this post with a question and perhaps the possible explanation:

As I write above, it took me a while to figure out how to wrap ``http.ResponseWriter`` correctly so that I could get access
to the HTTP status that was being set. The solution that was discussed in `this post <http://grokbase.com/t/gg/golang-nuts/12art4wedc/go-nuts-how-do-i-get-http-status-from-my-own-servehttp-function>`__ to just implement the ``WriteHeader()`` method didn't work for me.
``WriteHeader()`` method implemented by my ``MyResponseWriter()`` was never called except for then there was a redirect. I expected that
the call to ``Write()`` method of ``http.ResponseWriter()`` would invoke the version of ``WriterHeader()`` I implemented, but I cannot
see any way that could happen from the code in ``net/http/server.go``. So I think this is what's "implied" in this and all the other posts I have seen: the handler for the request must call ``WriteHeader()`` with the HTTP status as the server code above does.

It looks like `soon <https://github.com/golang/go/issues/18997>`__ there will be a direct way to get the HTTP response status.


References
==========

The following links helped me understand the above and write this post:

- http://jordanorelli.com/post/42369331748/function-types-in-go-golang
- https://golang.org/doc/effective_go.html#interface_methods
- https://gocodecloud.com/blog/2016/11/15/simple-golang-http-request-context-example/
- https://www.slideshare.net/blinkingsquirrel/customising-your-own-web-framework-in-go






