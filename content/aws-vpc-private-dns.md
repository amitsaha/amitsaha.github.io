Title: AWS Private Route53 DNS and Docker containers
Date: 2018-08-14 20:00
Category: infrastructure
Status: Draft

```
#!/bin/bash -ex
# Add VPC DNS server as an additional DNS server
# adapted from https://stackoverflow.com/questions/39100395/getting-the-dns-ip-used-within-an-aws-vpc
# and https://development.robinwinslow.uk/2016/06/23/fix-docker-networking-dns/
mac=$(curl -s -S http://169.254.169.254/latest/meta-data/mac)
cidr=$(curl -s -S "http://169.254.169.254/latest/meta-data/network/interfaces/macs/$mac/vpc-ipv4-cidr-block")
network=$(echo $cidr | cut -f 1 -d '/')
vpc_dns="$(echo $network  | cut -d'.' -f1-2 ).0.2"
echo '{"dns":["'"$vpc_dns"'", "8.8.8.8"]}' > /etc/docker/daemon.json
cat /etc/docker/daemon.json
systemctl restart docker
```
