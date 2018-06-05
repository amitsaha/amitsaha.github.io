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
 
