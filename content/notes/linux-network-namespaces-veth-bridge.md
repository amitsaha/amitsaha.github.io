Title: Linux Network Namespaces, veth pair and bridges
Date: 2018-08-14 10:00


Test bed:

```
ubuntu@ip-172-34-59-184:~$ sudo docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
efe90e4a2202        ubuntu              "bash"              About an hour ago   Up About an hour                        affectionate_nightingale
fef1fd3fde6c        ubuntu              "bash"              3 hours ago         Up 3 hours                              brave_ramanujan
ubuntu@ip-172-34-59-184:~$
```

Experiments:

NAT-ed traffic in default bridge mode:

```
$ sudo tcpdump -i eth0 port 10001
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on eth0, link-type EN10MB (Ethernet), capture size 262144 bytes
23:35:50.462020 IP ip-172-34-59-184.51116 > 172.34.1.252.10001: Flags [S], seq 1535892551, win 29200, options [mss 1460,sackOK,TS val 2391397203 ecr 0,nop,wscale 7], length 0
23:35:50.463033 IP 172.34.1.252.10001 > ip-172-34-59-184.51116: Flags [S.], seq 3768637232, ack 1535892552, win 64000, options [mss 8961,nop,wscale 0,sackOK,TS val 126715594 ecr 2391397203], length 0
23:35:50.463072 IP ip-172-34-59-184.51116 > 172.34.1.252.10001: Flags [.], ack 1, win 229, options [nop,nop,TS val 2391397204 ecr 126715594], length 0
```

veth/link/bridge configuration:

```
ubuntu@ip-172-34-59-184:~$ ip link show type veth
5: veth2d85589@if4: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue master docker0 state UP mode DEFAULT group default
    link/ether ba:0d:ce:ed:72:f8 brd ff:ff:ff:ff:ff:ff link-netnsid 0
7: veth6292b12@if6: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue master docker0 state UP mode DEFAULT group default
    link/ether 56:53:37:71:cc:2b brd ff:ff:ff:ff:ff:ff link-netnsid 1
ubuntu@ip-172-34-59-184:~$ ip link show type bridge
3: docker0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP mode DEFAULT group default
    link/ether 02:42:b6:1f:1e:51 brd ff:ff:ff:ff:ff:ff
ubuntu@ip-172-34-59-184:~$ sudo ip netns exec brave_ramanujan ip link show type veth
4: eth0@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP mode DEFAULT group default
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
ubuntu@ip-172-34-59-184:~$ sudo ip netns exec brave_ramanujan ip link show type bridge
ubuntu@ip-172-34-59-184:~$
```

## References

- https://blog.scottlowe.org/2013/09/04/introducing-linux-network-namespaces/
- https://platform9.com/blog/container-namespaces-deep-dive-container-networking/
- https://github.com/micahculpepper/dockerveth
- https://wiki.archlinux.org/index.php/Network_bridge
- http://man7.org/linux/man-pages/man4/veth.4.html
