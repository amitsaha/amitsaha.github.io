Title: Getting a chance before systemd shuts your supervisord down
Date: 2018-01-12 12:00
Category: fedora
Status: Draft

If you are running your server applications via [supervisord]() on a base distro running [systemd](), you may find 
this post useful.

## Scenario

An example scenario to help us establish the utility for this post is as followws:

- `systemd` starts the shutdown process
- `systemd` stops `supervisord`
- `supervisord` stops your processes
- You see in-flight requests being dropped

## Proposed solution

What we want to do is **prevent** in-flight requests being dropped. We can do so in two ways:

1. Our server application is intelligent enough to not exit without ever exiting if a request is in progress
2. We hook into the shutdown process above so that we stop new requests from coming in once the shutdown process has started and give our application server enough time to finish doing what it is doing.

The first appraoch has more theoretical "guarantee" around what we want, but can be hard to implement correctly. 
What if your clients keep long lived connections? You can never know when it might start sending you a new request. What
if you just can't do it right even after you have tried all sorts of signal handling tricks? 

The second approach is a less-pure way of doing it but it gives us a practical guarantee:

- We signal that we are shutting down (say enable maintenance mode in consul) a specific service instance
- Our proxy/load balancer detects the above event
- The proxy/load balancer stops sending traffic

When you are using a 


- essentially buy ourself sometime before `systemd` shuts the
`supervisord` service down. The proposed solution is a systemd unit - let's call it `drain-connections` which works
as follows:

- 




```
[Unit]
Description=Drain connections
After=supervisord.service
BindsTo=supervisord.service
Conflicts=shutdown.target reboot.target halt.target

[Service]
Type=oneshot
RemainAfterExit=true
ExecStart=/bin/true
ExecStop=/usr/local/bin/supervisorctl stop sheldon
ExecStop=/bin/sleep 300
TimeoutSec=300

[Install]
WantedBy=multi-user.target
```

```
[Unit]
Description=Process Monitoring and Control Daemon
After=rc-local.service
BindsTo=drainconnections
Before=drainconnections.service

[Service]
Type=forking
ExecStart=/usr/bin/supervisord -c /etc/supervisord.conf

[Install]
WantedBy=multi-user.target

```
