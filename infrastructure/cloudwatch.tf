resource "aws_cloudwatch_log_group" "log_group" {
  name              = "order-book-${var.environment_name}"
  retention_in_days = 7
  tags              = local.tags
}
