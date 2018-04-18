Title: My own ping
Date: 2018-02-19 07:40
Category: software
Status: Draft

## C


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
