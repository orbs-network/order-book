output "app_url" {
  value       = aws_apprunner_service.order-book.service_url
  description = "The URL of the deployed application"
}
