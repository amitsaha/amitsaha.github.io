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

What we want to do is prevent the last scenario - essentially buy ourself sometime before `systemd` shuts the
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
