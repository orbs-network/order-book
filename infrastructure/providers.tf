terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    heroku = {
      source  = "heroku/heroku"
      version = "~> 5.0"
    }
  }


  backend "s3" {}
}

provider "aws" {
  assume_role {
    role_arn     = "arn:aws:iam::506367651493:role/terraform"
    session_name = "Terraform"
  }
}
