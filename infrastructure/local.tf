data "aws_caller_identity" "current" {}

data "aws_availability_zones" "available" {}

locals {
  ecr_image = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.region}.amazonaws.com/perps-repo-${var.environment_name}:${var.image_tag}"

  tags = {
    account     = data.aws_caller_identity.current.account_id
    environment = var.environment_name
    region      = var.region
    service     = "order-book"
    owner       = "luke"
  }
}
