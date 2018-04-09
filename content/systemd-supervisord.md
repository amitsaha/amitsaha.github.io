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

## Solution

What we want to do is **prevent** in-flight requests being dropped when a system is shutting down as part of
a power off cycle (AWS instance termination, for example). We can do so in two ways:

1. Our server application is intelligent enough to not exit (and hence halt instance shutdown) if a request is in progress
2. We hook into the shutdown process above so that we stop new requests from coming in once the shutdown process has started and give our application server enough time to finish doing what it is doing.

The first approach has more theoretical "guarantee" around what we want, but can be hard to implement correctly. In fact,
I couldn't get it right even after trying all sorts of signal handling tricks. So, I went ahead with the very unclean
second approach:

- Register a shutdown "hook" which gets invoked when `systemd` wants to stop `supervisord`
- This hook takes the service instance out of the healthy pool
- The proxy/load balancer detects the above event and stops sending traffic
- As part of the "hook", after we have gotten ourself out of the healthy service pool, we sleep for an arbitary time so that
existing requests can finish

When you are using a software like [linkerd]() as your RPC proxy, even long-lived connections are not a problem since
`linkerd` will see that your service instance is unhealthy, so it will not proxy any more requests to it.


## Proposed solution implementation

The proposed solution is a systemd unit - let's call it `drain-connections` which is defined as follows:

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
