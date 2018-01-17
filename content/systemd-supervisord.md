Title: Setting up OpenVPN client with systemd template unit files
Date: 2018-01-12 12:00
Category: fedora
Status: Draft

```
[Unit]
Description=Drain Connections
After=supervisord
BindsTo=supervisord.service

[Service]
Type=oneshot
RemainAfterExit=True
ExecStart=/bin/true
ExecStop=/usr/bin/touch /var/shuttingdown1
ExecStop=/usr/bin/sleep 30
ExecStop=/usr/bin/touch /var/shuttingdown2

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
