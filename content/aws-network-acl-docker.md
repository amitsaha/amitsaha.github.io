Title: AWS Network ACLs, docker containers and ephermal port range
Date: 2018-07-26 16:00
Category: infrastructure
Status: Draft

Working with AWS VPC has given and continutes to give me a great opportunity to learn, revisit and relearn a lot of 
computer networking theory. In this post, I share something that involves AWS Network ACLs, docker containers and ephermal 
port ranges. I will try and put a bit of background in each section mainly so that it clears up things for me and may be
others who may need a refresher.

## Ephermal ports

Communication over network sockets involves two parties - usually referred to as a client and a server.

## Ephermal ports range

## Docker networks

## NAT on the host

## AWS Network ACLs



## Solutions

### Set the ephermal port range on the host

### Set the ephermal port range in a Linux docker container

https://docs.docker.com/engine/reference/commandline/run/#configure-namespaced-kernel-parameters-sysctls-at-runtime

###  Set the ephermal port range in a Windows docker container

https://support.microsoft.com/en-au/help/929851/the-default-dynamic-port-range-for-tcp-ip-has-changed-in-windows-vista
https://forums.docker.com/t/able-to-set-ephermal-port-range-in-windows-containers-via-docker-run/56095

```
PS C:\> netsh int ipv4 get dynamicport tcp
The following command was not found: int ipv4 get dynamicport tcp.
PS C:\> netsh int ipv4 show dynamicport tcp

Protocol tcp Dynamic Port Range
---------------------------------
Start Port      : 49152
Number of Ports : 16384

PS C:\> netsh int ipv4 set dynamicport tcp start=50000 num=1000
Ok.

PS C:\> netsh int ipv4 show dynamicport tcp

Protocol tcp Dynamic Port Range
---------------------------------
Start Port      : 50000
Number of Ports : 1000
```


