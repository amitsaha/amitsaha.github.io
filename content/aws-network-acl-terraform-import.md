Title: Importing existing AWS Network ACL into Terraform
Date: 2018-10-15 20:00
Category: infrastructure
Status: Draft

Recently, I worked on importing some AWS resources into Terraform. After importing a few
VPC resources - vpc, subnets, routes, routing tables, came the turn of [Network ACLs](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-network-acls.html). 

While configuring network ACL rules, we basically have two choices with Terraform 
(similar to security group and security group rules):

1. We can define the entries/rules inline in the `aws_network_acl` resource
2. We can define the entries/rules as separate `aws_network_acl_rule` resources

My first thought was I wanted to define separate `aws_network_acl_rule` resources. However,
I then learned that whereas Network ACLs could be imported, Network ACL entries
are not yet [supported](https://github.com/terraform-providers/terraform-provider-aws/issues/704#issuecomment-433181340).
Hence, I decided that I will define the network acl entries/rules inline so that importing the network acl would
also import the ACl entries/rules. 

However, the network ACLs I was trying to import had a fair number of ACL entries/rules in them and writing them by hand
would be a mind numbing exercise. The reasoning would have applied even if I were to write them as 
separate `aws_network_acl_rule` resources. So, I decided to generate the Terraform code using my hobby AWS CLI
tool - [yawsi](https://github.com/amitsaha/yawsi).



subnet association import

