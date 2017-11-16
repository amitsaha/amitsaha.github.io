Title: Demo of using consul tags
Date: 2017-11-16 18:00
Category: infrastructure


We have a server offering a service which is discovered via [consul](https://www.consul.io/). Consul tags
allow us to achieve different goals we may have in interacting with this service.

**Problem #1**

Let's say our service is a HTTP server acting as a routing point for multiple independent resources

**Solution**

When using `docker run`:

```
$ sudo docker run --add-host myhost.name:127.0.0.1 -ti python bash
Unable to find image 'python:latest' locally
latest: Pulling from library/python
Digest: sha256:eb20fd0c13d2c57fb602572f27f05f7f1e87f606045175c108a7da1af967313e
Status: Downloaded newer image for python:latest
...
```

This will show up as an additional entry in the container's `/etc/hosts` file:

```
root@fee9aeccbc4b:/# cat /etc/hosts
...
127.0.0.1	myhost.name
```

With `docker compose`, we can use the `extra_hosts` key:

```
extra_hosts:
    - "myhost.name:127.0.0.1"
```
