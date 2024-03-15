module "vpc" {
  source  = "cloudposse/vpc/aws"
  version = "2.1.0"

  name                    = "order-book-${var.environment_name}-vpc"
  ipv4_primary_cidr_block = var.vpc_cidr_block

  tags = local.tags
}

module "subnets" {
  source  = "cloudposse/dynamic-subnets/aws"
  version = "2.4.1"

  name                 = "order-book-${var.environment_name}"
  availability_zones   = local.azs
  vpc_id               = module.vpc.vpc_id
  igw_id               = [module.vpc.igw_id]
  ipv4_cidr_block      = [module.vpc.vpc_cidr_block]
  nat_gateway_enabled  = true
  nat_instance_enabled = true

  tags = local.tags
}
