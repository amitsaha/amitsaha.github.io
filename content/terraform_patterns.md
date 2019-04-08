Title: Creating multiple instances of a resource in Terraform
Date: 2019-04-04 16:00
Category: infrastructure
Status: Draft

In this post, we will use [Golang](https://golang.org/) to generate Terraform configuration from a TOML specification.

## Background on `count`

Using count is a [popular approach](https://www.terraform.io/docs/configuration/resources.html#count-multiple-resource-instances) to
creating multiple instances of the same resource. I have been combining it with `lists` and `maps` to configure
multiple instances of resources such as AWS VPC subnets, Autoscaling groups and most recently Network ACL rules. 

For example:

```

module "vpc_services" {
  source   = "../../modules/vpc"
  ...
  
  private_subnet_nacl_rules = "${list(
    map(
      "subnet_name", "SubnetA",
      "rule_number", 100,
      "egress", false,
      "protocol", "tcp",
      "rule_action", "allow",
      "cidr_block","${local.vpc_root}.12.0/24",
      "from_port", 1433,
      "to_port", 1433
    ),
    map(
      "name", "SubnetB",
      "rule_number", 101,
      "egress", true,
      "protocol", "tcp",
      "rule_action", "allow",
      "cidr_block","${local.vpc_root}.93.0/24",
      "from_port", 32768,
      "to_port", 65535
    ),
    ...
    # more such rules
  )}"

}

```

The resource creation looks as follows:

```
locals {
    public_network_acl_ids_map = "${zipmap(
        aws_subnet.public.*.tags.Name, aws_network_acl.public_subnets.*.id
    )}",
    private_network_acl_ids_map = "${zipmap(
        aws_subnet.private.*.tags.Name, aws_network_acl.private_subnets.*.id
    )}"
}

...

resource "aws_network_acl_rule" "private_subnet_rules" {
    count = "${length(var.private_subnet_nacl_rules)}"

    network_acl_id = "${lookup(
        local.private_network_acl_ids_map,
        lookup(var.private_subnet_nacl_rules[count.index], "subnet_name")
    )}"

    rule_number    = "${lookup(var.private_subnet_nacl_rules[count.index], "rule_number")}"
    egress         = "${lookup(var.private_subnet_nacl_rules[count.index], "egress")}"
    protocol       = "${lookup(var.private_subnet_nacl_rules[count.index], "protocol")}"
    rule_action    = "${lookup(var.private_subnet_nacl_rules[count.index], "rule_action")}"

    cidr_block     = "${lookup(var.private_subnet_nacl_rules[count.index], "cidr_block")}"
    from_port      = "${lookup(var.private_subnet_nacl_rules[count.index], "from_port")}"
    to_port        = "${lookup(var.private_subnet_nacl_rules[count.index], "to_port")}"
}
```

Since we are using the `count` attribute which Terraform uses in its state to keep track of the resources' state, 
a change in an item somewhere in the middle of the `private_subnet_nacl_rules` list, will in this case cause the 
rules following itto be created and destroyed. Of course, this is not limited to Network ACL rules. See [issue](https://github.com/hashicorp/terraform/issues/14275). 

What do we do? The most straightforward approach to this is to create separate `aws_network_acl_rule` resources
by hand. Instead of writing by hand however, what if we generate the ACL rules? That way:

- We don't run into the issue with count
- We don't have to manually write the terraform configuration for each network ACL rule

## Specification for Network ACL rules

An AWS network ACL rule has the following specification:

- Rule number
- Egress or ingress
- protocol
- from port
- to port
- CIDR block
- Network ACL to which it is attached to

I propose a [toml](https://github.com/toml-lang/toml) based specification:

```
subnet_name = "SubnetA"

rules = [
    {rule_no=101, egress = false, protocol = "tcp", rule_action = "allow", cidr_block = "127.0.0.1/32", from_port = 22, to_port = 30},
    {rule_no=102, egress = false, protocol = "tcp", rule_action = "allow", cidr_block = "127.0.0.1/32", from_port = 22, to_port = 30}
]
```
The assumption here is that, we will have a Network ACL rules specification file per Network ACL and the network ACL ID 
will be derived from the Subnet's name specified in `subnet_name`.

## Generating Terraform configuration

Now that we have a specification for our network acl rules, we will now write our program which will generate Terraform code 
from it. I will be using [burntsushi/toml](https://github.com/BurntSushi/toml) to parse the TOML file and serialize
it into a Golang structure. The key bit here is the Golang struct:

```
type naclRulesSpec struct {
	SubnetName string `toml:"subnet_name"`
	Rules      []naclRule
}

type naclRule struct {
	NetworkACLID string `tf:"network_acl_id"`
	Egress       bool   `toml:"egress" tf:"egress" tf_type:"bool"`
	RuleNo       int64  `toml:"rule_no" tf:"rule_number" tf_type:"int"`
	RuleAction   string `toml:"rule_action" tf:"rule_action"`
	CidrBlock    string `toml:"cidr_block" tf:"cidr_block"`
	Protocol     string `toml:"protocol" tf:"protocol"`
	FromPort     int64  `toml:"from_port" tf:"from_port" tf_type:"int"`
	ToPort       int64  `toml:"to_port" tf:"to_port" tf_type:"int"`
}

```











https://github.com/hashicorp/terraform/issues/17144
