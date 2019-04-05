Title: Creating multiple instances of the same resource in Terraform
Date: 2019-04-04 16:00
Category: infrastructure
Status: Draft

Using count is a [popular approach](https://www.terraform.io/docs/configuration/resources.html#count-multiple-resource-instances) to
creating multiple instances of the same resource. I have recently beeen combining with `lists` and `maps` to configure
multiple instances of resources such as AWS VPC subnets, Autoscaling groups and most recently Network ACL rules. 

For example:

```

module "vpc_services" {
  source   = "../../modules/vpc_services"
  ...
  
  private_subnet_nacl_rules = "${list(
    map(
      "subnet_name", "DatabaseA",
      "rule_number", 100,
      "egress", false,
      "protocol", "tcp",
      "rule_action", "allow",
      "cidr_block","${local.vpc_root}.98.0/24",
      "from_port", 1433,
      "to_port", 1433
    ),
    map(
      "name", "DatabaseB",
      "rule_number", 101,
      "egress", true,
      "protocol", "tcp",
      "rule_action", "allow",
      "cidr_block","${local.vpc_root}.97.0/24",
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

There's a problem with this which we run into especially with Network ACL rules which may be prone to frequent change. Since
we are using the `count` attribute which Terraform uses in its state to keep track of the resources' state, a change in an
item somewhere in the middle of the `private_subnet_nacl_rules` list, will in this case cause the rules following it
to be created and destroyed. Of course, this is not limited to Network ACL rules. See [issue](https://github.com/hashicorp/terraform/issues/14275).

What do we do? The most straightforward approach to this is to create separate `aws_network_acl_rule` resources.
However, the downside to it is to having to manually write these resources. Instead, let's generate the Terraform
configuration from a specification of the network ACL rules. That way:

- We don't run into the issue with count
- We don't have to manually write the terraform configuration for each network ACL rule

https://github.com/hashicorp/terraform/issues/17144
