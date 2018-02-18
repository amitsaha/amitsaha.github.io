Title: Quickstart: Using Sorted Sets in Redis from CLI, Python and Golang
Date: 2018-02-18 15:00
Category: software
Status: draft

In this post, I will first install a [redis](https://redis.io/) server on Fedora, and demo 
[redis sorted sets](https://redis.io/topics/data-types-intro) and talk to the server using `redis-cli` and clients
from Python and Golang. If you are running another operating system, please see the [download page](https://redis.io/download).

## Installation and server setup

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

## Check if Redis is alive

Once the server has started, let's check if our server is up and running:

```
$ redis-cli ping
PONG
```

## Sorted Sets

Redis' sorted set is a set data structure but each element in the set is also associated with a `score`. It is a
hash map but with an interesting property - the set is ordered based on this `score` value. 

This allows us to perform the following operations easily:

- Retrieve the top or bottom 10 keys based on the score
- Find the rank/position of a key in the set
- The score of a key can be updated anytime while the set will be adjusted (if needed) based on the new score

The section on sorted sets [here](https://redis.io/topics/data-types#sorted-sets) and [here](https://redis.io/topics/data-types-intro) in the Redis docs has more details on this.

## Example scenario: Top tags

We will now create a sorted set called `tags`. This set will store tags for posts in a blog or some other content
system where entries can have one or more tags associated with them. At any given point of time, we would like to
know what are the top 5 tags in our system.

## Implementation using `redis-cli`

We will first add a few tags to our sorted set `tags` using the [ZADD](https://redis.io/commands/zadd) command:

```
127.0.0.1:6379> ZADD tags 1 "python"
(integer) 1
127.0.0.1:6379> ZADD tags 1 "golang"
(integer) 1
127.0.0.1:6379> ZADD tags 1 "redis"
(integer) 1
127.0.0.1:6379> ZADD tags 1 "flask"
(integer) 1
127.0.0.1:6379> ZADD tags 1 "rust"
(integer) 1
127.0.0.1:6379> ZADD tags 2 "rust"
(integer) 0
127.0.0.1:6379> ZADD tags 3 "python"
(integer) 0
127.0.0.1:6379> ZADD tags 1 "docker"
(integer) 1
127.0.0.1:6379> ZADD tags 1 "linux"
(integer) 1
127.0.0.1:6379> ZADD tags 1 "c"
(integer) 1
127.0.0.1:6379> ZADD tags 1 "software"
(integer) 1
1

```

Above, I used the command to update the score of `rust` and `python` twice to be 2 and 3 respectively. I could have used
[ZINCRBY](https://redis.io/commands/zincrby) as well. Now, I will list all the keys using the [zrange]() command:

```
127.0.0.1:6379> zrange tags 0 -1
 1) "c"
 2) "docker"
 3) "flask"
 4) "golang"
 5) "linux"
 6) "memcache"
 7) "redis"
 8) "software"
 9) "rust"
10) "python"
```

Note how the last two keys are `rust` and `python` - as they have the highest scores (2 and 3 respectively). The others are
sorted lexicographically. 

To reverse the order, we will use the [zrevrange]() command:

```
127.0.0.1:6379> ZREVRANGE tags 0 -1 withscores
 1) "python"
 2) "3"
 3) "rust"
 4) "2"
 5) "redis"
 6) "1"
 7) "golang"
 8) "1"
 9) "flask"
10) "1"

```

Above, we can see how with the `withscores` command, we also get the scores back.

Now, to get the top 5 tags, we will do the following:

```
127.0.0.1:6379> ZREVRANGE tags 0 4 withscores
 1) "python"
 2) "3"
 3) "rust"
 4) "2"
 5) "software"
 6) "1"
 7) "redis"
 8) "1"
 9) "memcache"
10) "1"
1
```



## Resources

- [Redis data types](https://redis.io/topics/data-types-intro)
- [Redis quick start](https://redis.io/topics/quickstart)

