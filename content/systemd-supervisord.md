Title: Getting a chance before systemd shuts you down
Date: 2018-01-12 12:00
Category: fedora
Status: Draft

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
