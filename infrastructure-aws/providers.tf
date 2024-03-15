terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }


  backend "s3" {}
}

provider "aws" {
  region = var.region
  assume_role {
    role_arn     = "arn:aws:iam::${var.aws_deploy_account}:role/${var.aws_deploy_iam_role_name}"
    session_name = "Terraform"
  }
}
