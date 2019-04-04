Title: Few terraform patterns I have been using
Date: 2019-04-04 16:00
Category: infrastructure
Status: Draft

locals {
  vpc_root = "172.19"
}

module "vpc_services" {
  source   = "../../modules/vpc_services"
  vpc_name = "ProductionServices"
  vpc_root = "${local.vpc_root}"

  vpc_cidr   = "${local.vpc_root}.0.0/16"
  ip_address = "${module.static.ip_address}"

  environment = "production"

  private_subnet_configs = "${list(
    map(
      "name","DatabaseA",
      "availability_zone","ap-southeast-2a",
      "cidr_block","${local.vpc_root}.97.0/24",
      "route_table_name","DatabaseARouteTable",
    ),
    map(
      "name","DatabaseB",
      "availability_zone","ap-southeast-2c",
      "cidr_block","${local.vpc_root}.98.0/24",
      "route_table_name","DatabaseBRouteTable",
    ),
    map(
      "name","ReportingDbA",
      "availability_zone","ap-southeast-2a",
      "cidr_block","${local.vpc_root}.100.0/24",
      "route_table_name","ReportingDbRouteTable",
    ),
    map(
      "name","ReportingDbB",
      "availability_zone","ap-southeast-2c",
      "cidr_block","${local.vpc_root}.101.0/24",
      "route_table_name","ReportingDbRouteTable",
    ),
  )}"

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
      "subnet_name", "DatabaseA",
      "rule_number", 100,
      "egress", true,
      "protocol", "tcp",
      "rule_action", "allow",
      "cidr_block","${local.vpc_root}.98.0/24",
      "from_port", 32768,
      "to_port", 65535
    ),
    map(
      "name", "DatabaseB",
      "rule_number", 101,
      "egress", false,
      "protocol", "tcp",
      "rule_action", "allow",
      "cidr_block","${local.vpc_root}.97.0/24",
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
  )}"

  public_subnet_configs = "${list(
    map(
      "name","Management",
      "availability_zone","ap-southeast-2a",
      "cidr_block","${local.vpc_root}.65.0/24",
      "route_table_name","ManagementRouteTable",
    ),
   map(
      "name","NAT",
      "availability_zone","ap-southeast-2a",
      "cidr_block","${local.vpc_root}.66.0/24",
      "route_table_name", "NATRouteTable",
    ),
    map(
      "name","ReportingStaff",
      "availability_zone","ap-southeast-2a",
      "cidr_block","${local.vpc_root}.67.0/24",
      "route_table_name", "ReportingStaffRouteTable",
    ),
  )}"

  nat_instances_config = "${list(
       map(
           "name", "NAT"
       )
   )}"
}

output "vpc_id" {
  value = "${module.vpc_services.vpc_id}"
}

output "private_subnet_ids_map" {
  value = "${module.vpc_services.private_subnet_ids_map}"
}

output "public_subnet_ids_map" {
  value = "${module.vpc_services.private_subnet_ids_map}"
}

output "private_route_table_ids_map" {
  value = "${module.vpc_services.private_route_table_ids_map}"
}

output "private_subnet_cidr_blocks" {
  value = "${module.vpc_services.private_subnet_cidr_blocks}"
}

module "vpc_peering_production_legacy" {
  source = "../../modules/vpc_peering"

  vpc_peer_requester_vpc_id = "${module.vpc_services.vpc_id}"
  vpc_peer_peer_vpc_id      = "${module.static.vpc_ids["legacy_production"]}"
  vpc_peer_peer_owner_id    = "${module.static.aws_account_ids["legacy_production"]}"

  vpc_peer_requester_route_table_ids = [
    "${lookup(module.vpc_services.private_route_table_ids_map, "DatabaseARouteTable")}",
    "${lookup(module.vpc_services.private_route_table_ids_map, "DatabaseBRouteTable")}",
  ]
}

output "vpc_services_peer_legacy_prod" {
  value = "${module.vpc_peering_production_legacy.vpc_peering_ids}"
}


resource "aws_subnet" "private" {
  count = "${length(var.private_subnet_configs)}"

  vpc_id = "${aws_vpc.main.id}"

  cidr_block = "${lookup(
    var.private_subnet_configs[count.index], "cidr_block")
  }"

  availability_zone = "${lookup(
    var.private_subnet_configs[count.index], "availability_zone")
  }"

  tags = "${merge(
    var.common_private_subnet_tags,
    map(
      "Name", "${lookup(
        var.private_subnet_configs[count.index], "name"
      )}"
    ),
  )}"
}

output "private_subnet_ids" {
  value = "${aws_subnet.private.*.id}"
}

output "private_subnet_ids_map" {
  value = "${zipmap(
    aws_subnet.private.*.tags.Name, aws_subnet.private.*.id
  )}"
}

locals {
    public_network_acl_ids_map = "${zipmap(
        aws_subnet.public.*.tags.Name, aws_network_acl.public_subnets.*.id
    )}",
    private_network_acl_ids_map = "${zipmap(
        aws_subnet.private.*.tags.Name, aws_network_acl.private_subnets.*.id
    )}"
}

resource "aws_network_acl" "public_subnets" {
    count = "${length(var.public_subnet_configs)}"
    vpc_id     = "${aws_vpc.main.id}"
    subnet_ids = [
        "${aws_subnet.public.*.id[count.index]}"
    ]
    tags = {
        Name = "${lookup(var.public_subnet_configs[count.index], "name")}"
    }
}

resource "aws_network_acl" "private_subnets" {
    count = "${length(var.private_subnet_configs)}"
    vpc_id     = "${aws_vpc.main.id}"
    subnet_ids = [
        "${aws_subnet.private.*.id[count.index]}"
    ]
}

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



https://github.com/hashicorp/terraform/issues/17144