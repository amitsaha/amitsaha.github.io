```
#test.py

while True:
    f = open('test.txt', 'w')
    f.write('Hello')
    f.close()
```

```
$ python test.py &
```

```
vagrant@default-centos-7-latest:~$ sudo pidstat -p 13419 -d 1
Linux 4.13.0-21-generic (default-centos-7-latest)       06/05/2018      _x86_64_        (1 CPU)

04:02:48 AM   UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command
04:02:49 AM  1000     13419      0.00  31845.83      0.00       0  python
04:02:50 AM  1000     13419      0.00  29121.65      0.00       0  python
^C
Average:     1000     13419      0.00  30476.68      0.00       0  python
```

```
vagrant@default-centos-7-latest:~$ sudo pidstat -p 13419 1
Linux 4.13.0-21-generic (default-centos-7-latest)       06/05/2018      _x86_64_        (1 CPU)

04:02:53 AM   UID       PID    %usr %system  %guest   %wait    %CPU   CPU  Command
04:02:54 AM  1000     13419    3.03   22.22    0.00    4.04   25.25     0  python
04:02:55 AM  1000     13419    3.19   23.40    0.00    3.19   26.60     0  python
04:02:56 AM  1000     13419    2.11   23.16    0.00    3.16   25.26     0  python
^C
Average:     1000     13419    2.78   22.92    0.00    3.47   25.69     -  python
```

```
$  sudo perf trace -p 13419
..
 1509.315 ( 0.157 ms): openat(dfd: CWD, filename: 0x8601e010, flags: CREAT|TRUNC|WRONLY, mode: IRUGO|IWUGO) = 3
 1509.545 ( 0.028 ms): fstat(fd: 3</home/vagrant/test.txt>, statbuf: 0x7ffecd7dc2f0          ) = 0
 1509.697 ( 0.017 ms): fstat(fd: 3</home/vagrant/test.txt>, statbuf: 0x7ffecd7dc1b0          ) = 0
 1509.792 ( 0.029 ms): write(fd: 3</home/vagrant/test.txt>, buf: 0x564785fcf800, count: 5    ) = 5
 1509.877 ( 0.009 ms): close(fd: 3</home/vagrant/test.txt>                                   ) = 0
 ```
 
 
 ```
 vagrant@default-centos-7-latest:~$ sudo cat /proc/13419/wchan
io_schedule
```

```
vagrant@default-centos-7-latest:~$ sudo cat /proc/13419/stack
[<ffffffffba6b3226>] io_schedule+0x16/0x40
[<ffffffffba7b79c4>] wait_on_page_bit+0xf4/0x130
[<ffffffffba7cdb5d>] truncate_inode_pages_range+0x56d/0x830
[<ffffffffba7cdee8>] truncate_pagecache+0x48/0x70
[<ffffffffba9079c8>] ext4_setattr+0x8f8/0x9e0
[<ffffffffba872ca5>] notify_change+0x2e5/0x430
[<ffffffffba84cf33>] do_truncate+0x73/0xc0
[<ffffffffba861f28>] path_openat+0xf88/0x1630
[<ffffffffba8638db>] do_filp_open+0x9b/0x110
[<ffffffffba84f63b>] do_sys_open+0x1bb/0x2b0
[<ffffffffba84f764>] SyS_openat+0x14/0x20
[<ffffffffbaf0e8bb>] entry_SYSCALL_64_fastpath+0x1e/0xa9
[<ffffffffffffffff>] 0xffffffffffffffff
```
