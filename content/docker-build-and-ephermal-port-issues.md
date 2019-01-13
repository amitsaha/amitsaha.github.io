Title: Ephermal source port ranges and docker build
Date: 2019-01-14 20:00
Category: infrastructure
Status: Draft

I have written [previously](https://echorand.me/aws-network-acls-and-ephermal-port-ranges.html) about how things get interesting with ephermal port ranges in a Windows and Linux environment and AWS network acls. Today’s post is related to the same topic but specifically relevant if you are building docker images in such an environment.

Let’s start with the Dockerfile:

```
# Build runtime image
FROM microsoft/dotnet:2.2-aspnetcore-runtime
RUN apt-get -y update
..
```

The instruction `RUN apt-get -y update`  will make network requests to download resources from the Internet over HTTP. This means it will select a certain source port to make these HTTP requests. However, in a controlled environment, we want to explicitly state the range of ephermal ports that should be use, else these requests will not succeed. 

Let's see how we can do that.

## Background

How does a docker build happen? Inside containers. What do we do if we want to configure the ephermal port range for these builder containers? We can’t seem to be able to run sysctl in this scenario. We could use `docker build --host` to share the host’s network namespace. And that will ensure that out host’s ephermal port range will be used. However, we also had user namespacing turned on in our setup since this is a [sensible](https://echorand.me/docker-userns-remap-and-system-users-on-linux.html) thing to do. However, we cannot use a user namespace while using the host network. So, what do we do?

## Solution

Could we have a iptables rule to perform a source port translation so that anything that is going out of our host always uses a source port from the specified port range? Generally speaking, we will need to perform a variation of Source NAT. However, we will only change the source port and leave the IP address alone. The following rule will do it:

```
$ sudo iptables -t nat -I POSTROUTING -p tcp -m tcp --sport 32768:61000 -j MASQUERADE --to-ports 49152-61000
```

Here's what the above rule does:

1. We are adding this rule to the `nat` table (`-t nat`) in the `POSTROUTING`(`-I POSTROUTING`) chain
2. We want this rule to be applied for `TCP` (`-p tcp`) packets which has a source port in the range 32768-61000 (`--sport 32768-61000`)
3. If a packet matches our rule, forward it to the `MASQUERADE` target (`-j MASQUERADE`)
4. Once in the `MASQUERADE` target, change the source port to be in the range 49152-61000 (`--to-ports 49152-61000`)

There's much to learn about [iptables](https://www.frozentux.net/iptables-tutorial/iptables-tutorial.html). The most relevant
ones to be familiar with are:

- [Chains](https://www.booleanworld.com/depth-guide-iptables-linux-firewall/#Chains)
- [Masquerade target](https://www.frozentux.net/iptables-tutorial/chunkyhtml/x4422.html)
- [Iptables matches](https://www.frozentux.net/iptables-tutorial/iptables-tutorial.html#MATCHES)
