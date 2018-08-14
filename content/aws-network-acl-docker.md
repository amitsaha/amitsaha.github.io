Title: AWS Network ACLs, docker containers and ephermal port range
Date: 2018-07-26 16:00
Category: infrastructure
Status: Draft

Working with AWS VPC has given and continutes to give me a great opportunity to learn, revisit and relearn a lot of 
computer networking theory. In this post, I share something that involves AWS Network ACLs, docker containers and ephermal 
port ranges. I will try and put a bit of background in each section mainly so that it clears up things for me and may be
others who may need a refresher.

# Ephermal ports

Communication over network sockets involves two parties - usually referred to as a client and a server.

# Ephermal ports range

# Docker networks

# NAT on the host

# AWS Network ACLs



# Solutions

## Set the ephermal port range on the host

## Set the ephermal port range in a Linux docker container

##  Set the ephermal port range in a Windows docker container

netsh int ipv4 set dynamicport tcp start=10000 num=1000

https://support.microsoft.com/en-au/help/929851/the-default-dynamic-port-range-for-tcp-ip-has-changed-in-windows-vista
