module "redis" {
  source  = "cloudposse/elasticache-redis/aws"
  version = "1.2.0"

  namespace                     = "order-book"
  stage                         = var.environment_name
  name                          = "redis"
  description                   = "Redis instance for order book"
  vpc_id                        = module.vpc.vpc_id
  availability_zones            = local.azs
  subnets                       = module.subnets.private_subnet_ids
  create_security_group         = false
  associated_security_group_ids = [aws_security_group.redis_sg.id]
  cluster_size                  = var.multi_az_enabled == true ? 3 : 1
  instance_type                 = var.multi_az_enabled == true ? "cache.r7g.large" : "cache.t4g.micro"
  multi_az_enabled              = var.multi_az_enabled
  automatic_failover_enabled    = var.multi_az_enabled
  engine_version                = "7.1"
  family                        = "redis7"
  at_rest_encryption_enabled    = true
  transit_encryption_enabled    = false

  parameter = [
    {
      name  = "notify-keyspace-events"
      value = "lK"
    }
  ]

  log_delivery_configuration = [
    {
      destination      = aws_cloudwatch_log_group.log_group.name
      destination_type = "cloudwatch-logs"
      log_format       = "json"
      log_type         = "engine-log"
    }
  ]

  tags = local.tags
}
