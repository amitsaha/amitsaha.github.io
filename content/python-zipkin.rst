:Title: Introducing distributed tracing in your Python application via Zipkin
:Date: 2017-03-27 11:00
:Category: Python

Distributed tracing is the idea of tracing a network request as it travels through your services, as it would be in a microservices based architecture. The primary reason you may want to do is to troubleshoot or monitor the latency of a request
as it travels through the different services.

In this post we will see a demo of how we can introduce distributed tracing into a Python network stack communicating via HTTP. 
We have a service ``demo`` which is a Flask application, which listens on ``/``. The handler for ``/`` calls another service ``service1`` via HTTP. We want to be able to see how much time a request spends in each service by introducing distributed tracing. Before we get to the code, let's talk briefly about a few concepts.

Distributed Tracing concepts
============================

Roughly, a call to an "external service" starts a `span`. We can have a `span` nested within another span in a tree like fashion. All the spans in the context of a single request would form a `trace`. 

Something like the following would perhaps explain it better in the context of our ``demo`` and ``service`` network application stack:

.. code::

                   <--------------------          Trace       ------------------------------------ >                                                               Start Root Span                        Start a nested span      
   External Request -> Demo HTTP app       --->          Service 1 HTTP app        --->          Process
   

The span that is started from the ``service1`` is designated as a child of the ``root span`` which was started from the ``demo`` application. In the context of Python, we can think of a span as a context manager and one context manager living within another context manager. And all these "contexts" together forming a trace.

From the above it is somewhat clear (or not) that, the start of each span initiates a "timer" which then on the request's way back (or end of the span) is used to calculate the time the span lasted for. So, we need to have some thing (or things) which has to:

- Emit these data
- Recieve these data 
- Allow us to collate them together and make it available to us for each trace or request. 

This brings us to our next section.

Zipkin
======

`zipkin <http://zipkin.io/>`__ is a distributed tracing system which gives us the last two of the above requirements. How we emit these data from our application (the first point above) is dependent on the language we have written the application in and the distributed tracing system we chose for the last two requirements. In our case, `py_zipkin <https://github.com/Yelp/py_zipkin>`__ solves our problem.

First, we will start ``zipkin`` with ``elasticsearch`` as the backend as ``docker containers``. So, you need to have ``docker`` installed. To get the data in ``elasticsearch`` persisted, we will first create a `data container <http://echorand.me/data-only-docker-containers.html>`__ as follows:

.. code:

    $ docker create --name esdata openzipkin/zipkin-elasticsearch
    
Then, download my code from `here <https://github.com/amitsaha/python-web-app-recipes/archive/zipkin_python_demo.zip>`__ and:

.. code::

    $ wget ..
    $ unzip ..
    $ cd tracing/http_collector
    $ ./start_zipkin.sh


 
 
 
    
    






