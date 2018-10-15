Title: AWS VPC subnets and Internet connectivity
Date: 2018-10-15 20:00
Category: infrastructure
Status: Draft

We can have two kinds of subnets inside a AWS VPC - __private__ and __public__. What is the difference between
the two when it comes to Internet connectivity? The only difference is in one having a specific routing 
table entry and the other not having it when it comes to outgoing "Internet" connections. When you 
have a routing table entry inside the subnet's routing table sending traffic to an **Internet Gateway**, 
it's a __public__ subnet. On the other hand, if traffic from a subnet destined for the "Internet" is sent to either a 
NAT instance, or a AWS managed NAT device, the subnet is a __private__ subnet.

Let's discuss the flow of traffic

## Flow of traffic

## VPC Flow logs


https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html

