Title: Brief overview of using consul tags
Date: 2017-11-16 18:00
Category: infrastructure

[consul](https://www.consul.io/) allows a service to associate itself with `tags`. These are arbitrary
metadata that can be associated with the service and can be used for different purposes. Below I outline
a few examples of making use of tags and discuss some related topics.

## Use case #1: Dedicated service instances based on requests

Let's say our service is a HTTP server (REST API) acting as a routing point for multiple 
independent resources with the following service definition:


```
{
  "service": {
    "name": "api",
    "address": "",
    "port": 8080,
    "checks": [],
  }
}
```

We can then communicate with service using the DNS, `api.service.consul`.

Let's assume we are running N copies of this service, but want to have dedicated sub-pools for 
separate resource groups. For demonstration, let's say we have 2 separate pools. We will assign
the services in each pool a different tag as follows:

**projects**

```
{
  "service": {
    "name": "api",
    "address": "",
    "port": 8080,
    "checks": [],
    "tags":["projects"],
  }
}
```


**users**

```
{
  "service": {
    "name": "api",
    "address": "",
    "port": 8080,
    "checks": [],
    "tags":["users"],
  }
}
```

Once we register the services using the different tags, they automatically become discoverable via DNS
as `projects.api.service.consul` and `users.api.service.consul` respectively. Assuming that the routing 
to our HTTP server is happening in a higher layer, we will then direct traffic to these pools as follows:

```
api/projects/ -> projects.api.service.consul
api/users/ -> users.api.service.consul
```

## Use case #2: Running different versions of your service

We can use tags to run two different versions of our application for testing, gathering
performance data or any other reason:

**v1**

```
{
  "service": {
    "name": "api",
    "address": "",
    "port": 8080,
    "checks": [],
    "tags":["v1"],
  }
}
```


**v2**

```
{
  "service": {
    "name": "api",
    "address": "",
    "port": 8080,
    "checks": [],
    "tags":["v2"],
  }
}
```

We can then use a tag based weighted mechanism at a higher level proxy (such as [linkerd](https://github.com/linkerd/linkerd/commit/718514fb1d4b86153820880162d3c9559e115725)) to send traffic to these different service
versions.

## Use case #3: Other metadata

This [issue](https://github.com/hashicorp/consul/issues/997/) on consul's project discusses using
tags as a way to have artbitary metadata for a service due to the lack of support for key-value
labels.


## Using tags for discovery

Besides using the DNS interface for communicating with the services, we can use `tags` as filter with
the consul [catalog API](https://www.consul.io/api/catalog.html). However, it currently supports a single
tag. There is a feature request [open](https://github.com/hashicorp/consul/issues/1781) to support multiple
tags.


## Demo: Running two versions of a service

I have two versions of a service, `api`. The `v1` of the service is running on port `8080` and `v2` running
on port `8081`. `v1` and `v2` are also the tags associated with the respective instances.

The demo source code can be found [here](https://github.com/amitsaha/

DNS

Query API
