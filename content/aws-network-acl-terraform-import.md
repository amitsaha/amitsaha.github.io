Title: Importing existing AWS Network ACL into Terraform
Date: 2018-10-15 20:00
Category: infrastructure
Status: Draft

Recently, I worked on importing some AWS resources into Terraform. After importing a few
VPC resources - vpc, subnets, routes, routing tables, came the turn of Network ACLs. 
While configuring network ACL rules, we basically have two choices
with Terraform (similar to security group and security group rules):

1. We can define the entries/rules inline in the `aws_network_acl` resource
2. We can define the entries/rules as separate `aws_network_acl_rule` resources

My first thought was I wanted to define separate `aws_network_acl_rule` resources. However,
I then learned that where as Network ACLs could be imported, Network ACL entries
are not yet [supported](https://github.com/terraform-providers/terraform-provider-aws/issues/704#issuecomment-433181340).


The Network ACLs I was trying to import had a fair number of ACL entries/rules in them.


Hence, what I didn't want to do was manually write those ingress and agress rules either using the
inline format on the `aws_network_acl` resource or as top-level resources using the 
`aws_network_acl_rule` resource. 


The Network ACLs in question had a fair number of 

subnet association import

