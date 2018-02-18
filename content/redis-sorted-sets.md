Title: Quickstart: Using Sorted Sets in Redis from CLI, Python and Golang
Date: 2018-01-25 15:00
Category: software
Status: draft

In this post, I will first install a [redis](https://redis.io/) server on Fedora, and demo 
[redis sorted sets](https://redis.io/topics/data-types-intro) and talk to the server using `redis-cli` and clients
from Python and Golang. If you are running another operating system, please see the [download page](https://redis.io/download).

## Installation and setup server

We can install `redis` server using `dnf`, like so:

```
$ sudo dnf install redis
..
$ redis-server --version
Redis server v=4.0.6 sha=00000000:0 malloc=jemalloc-4.5.0 bits=64 build=427484a80e1b4515
```

Let's start the server:

```
$ sudo systemctl start redis
$ sudo systemctl status redis
● redis.service - Redis persistent key-value database
   Loaded: loaded (/usr/lib/systemd/system/redis.service; disabled; vendor preset: disabled)
  Drop-In: /etc/systemd/system/redis.service.d
           └─limit.conf
   Active: active (running) since Sun 2018-02-18 00:08:22 AEDT; 4s ago
 Main PID: 29944 (redis-server)
    Tasks: 4 (limit: 4915)
   CGroup: /system.slice/redis.service
           └─29944 /usr/bin/redis-server 127.0.0.1:6379

Feb 18 00:08:22 fedora.home systemd[1]: Starting Redis persistent key-value database...
Feb 18 00:08:22 fedora.home systemd[1]: Started Redis persistent key-value database.
..
```


Once the server has started, let's check if our server is up and running:

```
$ redis-cli ping
PONG
```

## Sorted Sets

Redis' sorted set is a set data structure but each element in the set is also associated with a `score`. The set is
ordered based on this `score` value. This allows us to perform the following operations easily:

- Retrieve the top or bottom 10 keys based on the score
- Find the rank/position of a key in the set
- The score of a key can be updated anytime

The section on [sorted sets](https://redis.io/topics/data-types-intro) in the Redis docs has more details on this.


## Resources

- [Redis data types](https://redis.io/topics/data-types-intro)
- [Redis quick start](https://redis.io/topics/quickstart)

