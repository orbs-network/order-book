variable "aws_deploy_account" {
  description = "AWS account id to deploy to"
  type        = string
  default     = "506367651493"
}

variable "aws_deploy_iam_role_name" {
  description = "AWS IAM role name to assume for deployment"
  type        = string
  default     = "terraform"
}

variable "environment_name" {
  type        = string
  description = "Environment specific name"
}

variable "region" {
  type        = string
  default     = "ap-northeast-1"
  description = "AWS region used"
}

variable "image_tag" {
  type        = string
  description = "Docker image tag"
}

variable "vpc_cidr_block" {
  type        = string
  description = "The CIDR block for the VPC"
  default     = "10.1.0.0/16"
}

variable "multi_az_enabled" {
  type        = bool
  description = "Enable multi AZ redundancy"
  default     = false
}

variable "az_count" {
  type        = number
  description = "Number of availability zones to use"
  default     = 3
}
