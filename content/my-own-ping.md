Title: How does ping roughly work over IPv4 on Linux?
Date: 2018-02-19 07:40
Category: software
Status: Draft

The `ping` program is one of the most common programs which is used to check the "aliveness" of a host and
a typical execution looks as follows:

```
$ ping 127.0.0.1 -c 1 -4

PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.062 ms

--- 127.0.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.062/0.062/0.062/0.000 ms
```

The `-c` switch indicates that we want to send a single "probe". The `-4` switch limits the `ping` program to stay
confined to making network operations related to IPv4 only.

It basically works by sending a special network packet to your destination host and waits for the host to 
reply back. Then, it prints if any packets were lost and the timing statistics. I wanted to understand 
how the program works - what does it send? what does it receive? The final product ideally would be a 
C program which will be a basic version of `ping`.

## Theory

This [pdf](http://www.galaxyvisions.com/pdf/white-papers/How_does_Ping_Work_Style_1_GV.pdf) here has a good description
of the worksing of ping. The non-detailed version is that we create a special ICMP packet, package it up within a IP 
packet and send it across to the destination. The destination Linux kernel receives the packet, and sends a reply
ICMP packet embedded within a IP packet. The destination host doesn't have any user space program running to receive
the "ping" packet. Each packet only has `header` information. You can embed specific data into the ICMP packet, but
that's not required. The post [here](http://www.genetech.com.au/blog/?p=970) describes the packet structure a bit 
more along with a graphical representation.

With that bit of theory under our belt, let's look into what system calls are made as part of the above invocation
of the `ping` program using `strace`.

## System calls made as part of ping

If you don't have `strace` installed, please install it using your package manager. Let's now execute the above `ping` program under `strace`:

```
$ strace -f ping 127.0.0.1 -c 1 -4
```

You will see the output of the above command similar to:

```
$ strace -e trace=network ping 127.0.0.1 -c 1 -4
socket(AF_INET, SOCK_DGRAM, IPPROTO_ICMP) = 3
socket(AF_INET, SOCK_DGRAM, IPPROTO_IP) = 4
connect(4, {sa_family=AF_INET, sin_port=htons(1025), sin_addr=inet_addr("127.0.0.1")}, 16) = 0
getsockname(4, {sa_family=AF_INET, sin_port=htons(34117), sin_addr=inet_addr("127.0.0.1")}, [16]) = 0
setsockopt(3, SOL_IP, IP_RECVERR, [1], 4) = 0
setsockopt(3, SOL_IP, IP_RECVTTL, [1], 4) = 0
setsockopt(3, SOL_IP, IP_RETOPTS, [1], 4) = 0
setsockopt(3, SOL_SOCKET, SO_SNDBUF, [324], 4) = 0
setsockopt(3, SOL_SOCKET, SO_RCVBUF, [65536], 4) = 0
getsockopt(3, SOL_SOCKET, SO_RCVBUF, [131072], [4]) = 0
PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
setsockopt(3, SOL_SOCKET, SO_TIMESTAMP, [1], 4) = 0
setsockopt(3, SOL_SOCKET, SO_SNDTIMEO, "\1\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 16) = 0
setsockopt(3, SOL_SOCKET, SO_RCVTIMEO, "\1\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 16) = 0
sendto(3, "\10\0q9\0\0\0\1\254k\331Z\0\0\0\0B,\0\0\0\0\0\0\20\21\22\23\24\25\26\27"..., 64, 0, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("127.0.0.1")}, 16) = 64
recvmsg(3, {msg_name={sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("127.0.0.1")}, msg_namelen=128->16, msg_iov=[{iov_base="\0\0x\314\0m\0\1\254k\331Z\0\0\0\0B,\0\0\0\0\0\0\20\21\22\23\24\25\26\27"..., iov_len=192}], msg_iovlen=1, msg_control=[{cmsg_len=32, cmsg_level=SOL_SOCKET, cmsg_type=0x1d /* SCM_??? */}, {cmsg_len=20, cmsg_level=SOL_IP, cmsg_type=IP_TTL, cmsg_data=[64]}], msg_controllen=56, msg_flags=0}, 0) = 64
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.188 ms

--- 127.0.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.188/0.188/0.188/0.000 ms
+++ exited with 0 +++
a
```

Let's look at the first three lines of the trace:

_socket(AF_INET, SOCK_DGRAM, IPPROTO_ICMP) = 3_

The above creates a socket of type `SOCK_DGRAM` and the protocol as `IPPROTO_ICMP`. The IPPROTO_ICMP socket
protocol was [added](https://lwn.net/Articles/443051/) to allow a friendlier way to create ICMP packets. This
eliminates the need to create "RAW" sockets which in turn eliminates the need to have the 
[CAP_NET_RAW](http://man7.org/linux/man-pages/man7/capabilities.7.html) capability. The file descriptor
returned is important to note here - `3`.

_socket(AF_INET, SOCK_DGRAM, IPPROTO_IP) = 4_

This creates another socket with `IPPROTO_IP` protocol and then uses it to connect to the UDP port 1025 on
the target host:

_connect(4, {sa_family=AF_INET, sin_port=htons(1025), sin_addr=inet_addr("127.0.0.1")}, 16) = 0_

The above two steps are needed to figure out the IP address of the network interface that will be used
to send the ICMP packets to the destination host. 

I am not quite sure why we need the new socket, hence I created an [issue] on the `iputils` project (https://github.com/iputils/iputils/issues/125) to request any clarification.

```
$ ./a.out
Error creating socket: Permission denied
```

```
$ sudo sysctl net.ipv4.ping_group_range
net.ipv4.ping_group_range = 1   0
```

```
$ id -g
1000
```

(If you are a member of multiple groups, the range has to include only one of the groups)

```
$ sudo sysctl -w net.ipv4.ping_group_range="0 1000"
net.ipv4.ping_group_range = 0 2000
asaha@asaha-desktop:~$ ./a.out
```

## Implementation

Next, we will look at the implementation of `ping` which is part of the [iputils](https://github.com/iputils/iputils) suite
of programs. However, before that, let's first see which network related system calls are made as part of the ping program
execution. 

https://lwn.net/Articles/443051/

### Resources

- [Stackoverflow](https://stackoverflow.com/questions/8290046/icmp-sockets-linux)
- ["ping" socket](https://lwn.net/Articles/443051/)
