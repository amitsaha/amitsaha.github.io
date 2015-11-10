A look at ss: another utility to investigate sockets

```

$ ss
Netid State      Recv-Q Send-Q                                                      Local Address:Port                                                          Peer Address:Port   
u_seq ESTAB      0      0                                                                  @00120 1260126                                                                  * 1260127
u_str ESTAB      0      0                                                    /var/run/docker.sock 1748357                                                                  * 1748356
u_str ESTAB      0      0                                                                       * 1453627                                                                  * 1451643
u_str ESTAB      0      0                                             /run/systemd/journal/stdout 27239                                                                    * 28320  
u_str ESTAB      0      0                                                      @/tmp/.X11-unix/X0 25666                                                                    * 25665  
u_str ESTAB      0      0                                                                       * 19166                                                                    * 19167  
u_str ESTAB      0      0                                                    /var/run/docker.sock 1757194                                                                  * 1756683

```
