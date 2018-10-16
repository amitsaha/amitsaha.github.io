Title: AWS VPC subnets and Internet connectivity over IPv4
Date: 2018-10-15 20:00
Category: infrastructure
Status: Draft

We can have two kinds of subnets inside a AWS VPC - __private__ and __public__. A public subnet is one which is 
attached to an [Internet Gateway](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html). This essentially adds a routing table entry to the subnet's routing table sending all Internet traffic to an **Internet Gateway**. On the other hand, 
if traffic from a subnet destined for the "Internet" is sent to either a NAT instance, or a AWS managed NAT device, the subnet 
is a __private__ subnet. 

An EC2 instance running in either subnet can choose to have a [public IP address](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-instance-addressing.html#concepts-public-addresses)
or not. Once you give your instance a public IP address, it becomes "reachable" from the Internet (let's call it `ingress`) 
and of course can reach "Internet" resources from the instance (let's call it `egress`).

Let's discuss the flow of traffic and the address translations that happens in various cases that may arise 
in the above scenario.

## Public subnet - Public IP

Consider a EC2 instance, E in a public subnet having a public IP.

### Ingress

Internet resource, device A talks to our instance E using the public IP. E sees it gets a request from A,
responds accordingly. The response goes via Internet Gateway configured for the subnet. The Internet Gateway
perfroms a Source Network Address translation where it changes the source IP address of the response packet
to match the public IP address of E, so that A gets the response from E, and not E's internal IP address.

### Egress

Instance E tries to access an Internet resource, B. The traffic goes via the Internet Gateway where a network 
address translation takes place - the source IP address is changed from the internal IP to the public facing
IP. When the response is received by the Internet Gateway, the destination IP is changed from the public IP of instance 
E to the private IP.

## Public subnet - No Public IP

Considering the above address translation, AWS will *not* prevent you from having an instance with no public IP,
but you will not be able to talk to any "Internet" resource and of course, vice-versa - i.e. no ingress or egress
from/to Internet resource(s).

## Private subnet - Public IP

Consider a EC2 instance, F in a private subnet having a public IP.

### Ingress

Internet resource, device B can talk to instance, F using the public IP. F gets the connection request, sends
back an acknowledgement. This acknowledgement needs to go to B via either a NAT instance or a AWS managed NAT device.
The NAT device performs a Source address translation replacing the private IP of the instance, F by the NAT's public
IP address. The acknowledgement makes it way back to the device B's "edge" device (firewall for example) and doesn't recognize
the public NAT IP as part of an ongoing connection and just drops it. Hence, although, device B can attempt to
initiate connection to F, it never actually succeeds.

## Egress

F tries to talk an Internet resource. The routing table forwards the request to the NAT instance/device. This device
then replaces the source IP of the request which currently is the private IP by the NAT device's Public IP addres and
forwards the request to the destination resource. When the response is received, the NAT device changes the destination
IP address from itself to that of the private IP address of the instance.


## Private subnet - No Public IP

Consider a EC2 instance, G in a private subnet having no public IP

## Ingress

In this case, there is no Ingress traffic to the instance, since there is no public IP. 

## Egress

G tries to talk an Internet resource. The routing table forwards the request to the NAT instance/device. This device
then replaces the source IP of the request which currently is the private IP by the NAT device's Public IP addres and
forwards the request to the destination resource. When the response is received, the NAT device changes the destination
IP address from itself to that of the private IP address of the instance.



## VPC Flow logs


https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html

