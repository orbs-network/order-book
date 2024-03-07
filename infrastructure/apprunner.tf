resource "aws_apprunner_auto_scaling_configuration_version" "order-book-asg" {
  auto_scaling_configuration_name = "order-book-asg-${var.environment_name}"
  min_size                        = 1
  max_size                        = 2
  tags                            = local.tags
}

resource "aws_apprunner_service" "order-book" {
  auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.order-book-asg.arn

  service_name = "order-book-${var.environment_name}"

  instance_configuration {
    instance_role_arn = aws_iam_role.apprunner_instance_role.arn
  }

  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_access_role.arn
    }

    image_repository {
      image_identifier      = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.region}.amazonaws.com/order-book-repo-${var.environment_name}:${var.image_tag}"
      image_repository_type = "ECR"

      image_configuration {
        runtime_environment_variables = {
          # TODO: update
          "REDIS_URL"           = "redis://${module.redis.endpoint}:${module.redis.port}"
          "RPC_URL"             = "http://rpc:8080"
          "PORT"                = "8080"
          "REPORT_SEC_INTERVAL" = "-1"
        }
      }
    }

    auto_deployments_enabled = false

  }

  tags = local.tags
}


resource "aws_apprunner_vpc_connector" "connector" {
  vpc_connector_name = "order-book-${var.environment_name}-connector"
  subnets            = module.subnets.public_subnet_ids
  security_groups    = [module.vpc.vpc_default_security_group_id]
}
