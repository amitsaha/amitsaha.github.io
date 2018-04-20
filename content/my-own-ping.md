Title: Roughly understanding how ping works over IPv4 on Linux
Date: 2018-02-19 07:40
Category: software
Status: Draft

The `ping` program is one of the most common programs which is used to check the "aliveness" of a host and
a typical execution looks as follows:

https://lwn.net/Articles/443051/



```
$ strace -f ping 127.0.0.1 -c 1 -4
```


```
socket(AF_INET, SOCK_DGRAM, IPPROTO_ICMP) = 3
capget({version=_LINUX_CAPABILITY_VERSION_3, pid=0}, NULL) = 0
capget({version=_LINUX_CAPABILITY_VERSION_3, pid=0}, {effective=0, permitted=0, inheritable=0}) = 0
socket(AF_INET, SOCK_DGRAM, IPPROTO_IP) = 4
connect(4, {sa_family=AF_INET, sin_port=htons(1025), sin_addr=inet_addr("127.0.0.1")}, 16) = 0
getsockname(4, {sa_family=AF_INET, sin_port=htons(43365), sin_addr=inet_addr("127.0.0.1")}, [16]) = 0
close(4)                                = 0
setsockopt(3, SOL_IP, IP_RECVERR, [1], 4) = 0
setsockopt(3, SOL_IP, IP_RECVTTL, [1], 4) = 0
setsockopt(3, SOL_IP, IP_RETOPTS, [1], 4) = 0
setsockopt(3, SOL_SOCKET, SO_SNDBUF, [324], 4) = 0
setsockopt(3, SOL_SOCKET, SO_RCVBUF, [65536], 4) = 0
getsockopt(3, SOL_SOCKET, SO_RCVBUF, [131072], [4]) = 0
fstat(1, {st_mode=S_IFCHR|0620, st_rdev=makedev(136, 2), ...}) = 0
write(1, "PING 127.0.0.1 (127.0.0.1) 56(84"..., 49PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
) = 49
setsockopt(3, SOL_SOCKET, SO_TIMESTAMP, [1], 4) = 0
setsockopt(3, SOL_SOCKET, SO_SNDTIMEO, "\1\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 16) = 0
setsockopt(3, SOL_SOCKET, SO_RCVTIMEO, "\1\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 16) = 0
rt_sigaction(SIGINT, {sa_handler=0x55e0c30fc440, sa_mask=[], sa_flags=SA_RESTORER|SA_INTERRUPT, sa_restorer=0x7fe01d753140}, NULL, 8) = 0
rt_sigaction(SIGALRM, {sa_handler=0x55e0c30fc440, sa_mask=[], sa_flags=SA_RESTORER|SA_INTERRUPT, sa_restorer=0x7fe01d753140}, NULL, 8) = 0
rt_sigaction(SIGQUIT, {sa_handler=0x55e0c30fc430, sa_mask=[], sa_flags=SA_RESTORER|SA_INTERRUPT, sa_restorer=0x7fe01d753140}, NULL, 8) = 0
rt_sigprocmask(SIG_SETMASK, [], NULL, 8) = 0
ioctl(1, TCGETS, {B38400 opost isig icanon echo ...}) = 0
ioctl(1, TIOCGWINSZ, {ws_row=14, ws_col=115, ws_xpixel=0, ws_ypixel=0}) = 0
sendto(3, "\10\0(l\0\0\0\1>,\330Z\0\0\0\0\3618\t\0\0\0\0\0\20\21\22\23\24\25\26\27"..., 64, 0, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("127.0.0.1")}, 16) = 64
setitimer(ITIMER_REAL, {it_interval={tv_sec=0, tv_usec=0}, it_value={tv_sec=10, tv_usec=0}}, NULL) = 0
recvmsg(3, {msg_name={sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("127.0.0.1")}, msg_namelen=128->16, msg_iov=[{iov_base="\0\0000\20\0\\\0\1>,\330Z\0\0\0\0\3618\t\0\0\0\0\0\20\21\22\23\24\25\26\27"..., iov_len=192}], msg_iovlen=1, msg_control=[{cmsg_len=32, cmsg_level=SOL_SOCKET, cmsg_type=0x1d /* SCM_??? */}, {cmsg_len=20, cmsg_level=SOL_IP, cmsg_type=IP_TTL, cmsg_data=[64]}], msg_controllen=56, msg_flags=0}, 0) = 64
write(1, "64 bytes from 127.0.0.1: icmp_se"..., 5764 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.059 ms

```

http://www.galaxyvisions.com/pdf/white-papers/How_does_Ping_Work_Style_1_GV.pdf

https://github.com/iputils/iputils/issues/125


http://www.genetech.com.au/blog/?p=970


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


### Resources

- [Stackoverflow](https://stackoverflow.com/questions/8290046/icmp-sockets-linux)
- ["ping" socket](https://lwn.net/Articles/443051/)
