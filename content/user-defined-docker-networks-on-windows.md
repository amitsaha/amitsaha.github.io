Title: User-defined networks in Docker for inter-container communication on Windows
Date: 2017-10-26 15:00
Category: infrastructure
Status: draft



On Windows 10, multiple `nat` networks are supported, but on Windows Server with 17.06 EE docker engine, only one `nat` network
is supported. Hence, when using `docker-compose`, we must specify the following:

```
networks:
  default:
    external:
      name: nat
      
```

The reason for the above is the default behavior of `docker-compose` is to create a new network for the services which will
fail with an error: `Problem : Error response from daemon: HNS failed with error : The parameter is incorrect`.
https://docs.microsoft.com/en-us/virtualization/windowscontainers/container-networking/network-drivers-topologies
