Title: Ephermal Port Ranges and Docker build
Date: 2018-08-14 20:00
Category: infrastructure
Status: Draft

I have written [previously](https://echorand.me/aws-network-acls-and-ephermal-port-ranges.html) about how things get interesting with ephermal port ranges in a Windows  and Linux environment, AWS network acls and ephermal port ranges. Today’s post is related to the same topic but specifically relevant if you are building docker images in such an environment.

Let’s start with the Dockerfile:

```
# Build runtime image
FROM microsoft/dotnet:2.2-aspnetcore-runtime
RUN apt-get -y update
..
```

The instruction `RUN apt-get -y update`  will make network requests to download resources from the Internet over HTTP. This means, it will select a certain source port to make these HTTP requests. 

How does a docker build happen? Inside containers. What do we do if we want to configure the ephermal port range for these builder containers? We can’t seem to be able to run sysctl in this scenario. We could use `docker build --host` to share the host’s network namespace. And that will ensure that out host’s ephermal port range will be used. However, we also had user namespacing turned on in our setup since this is a [sensible](https://echorand.me/docker-userns-remap-and-system-users-on-linux.html) thing to do. However, we cannot use a user namespace while using the host network. So, what do we do?

Enter iptables. Could we have a iptables rule to perform a source port translation so that anything that is going out of our host always uses a source port from the specified port range? Yes, the following rule will do it:

```
$ sudo iptables -t nat -I POSTROUTING -p tcp -m tcp --sport 32768:61000 -j MASQUERADE --to-ports 49152-61000
```
