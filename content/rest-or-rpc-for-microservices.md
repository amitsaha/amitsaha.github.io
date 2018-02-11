Title: Why RPC in Microservices?
Date: 2018-02-08 18:00
Category: software
Status: Draft

At Freelancer.com, my team works with different services. The [REST API](https://developers.freelancer.com/)
is powered by a number of backend services. The API itself is powered by a Python HTTP server which
communicates with the other services via RPC calls implemented using [Apache Thrift](https://thrift.apache.org/).
It is only during the past 2.5 years that I have been working with Apache Thrift or cross-language RPCs in general.
The question often comes up - why not just use HTTP throughout across all services? HTTP 1.1 is simple and easy
to understand. Implementing a HTTP API endpoint doesn't also mean having to learn about Apache Thrift, because
the data has to come from somewhere right. 

**Why not use HTTP throughout?**

Without going into benchmarks, the one reason I think RPCs are a better fit than HTTP is very well put in this
[blog post](https://blog.bugsnag.com/grpc-and-microservices-architecture/) on why they chose RPC:

> Unfortunately, it felt like we were trying to shoehorn simple methods calls into a data-driven RESTful interface. 
> The magical combination of verbs, headers, URL identifiers, resource URLs, and payloads that satisfied a RESTful 
> interface and made a clean, simple, functional interface seemed an impossible dream. RESTful has lots of rules and
> interpretations, which in most cases result in a RESTish interface, which takes extra time and effort to maintain its purity.


To implement a new interface, you first have to write the interface in some Interface Description Language (IDL).
You then have to generate the code in the language corresponding to your client(s) and server, then implement
the functionality itself. Changing an existing function signature is tricky business with Apache Thrift as well.
