variable "environment_name" {
  type        = string
  description = "Environment specific name"
}

variable "rpc_url" {
  type        = string
  description = "Blockchain RPC URL for the application"
  sensitive   = true
}

variable "log_level" {
  type        = string
  description = "Log level for the application"
  default     = "info"
}

variable "swap_contract_address" {
  type        = string
  description = "Swap reactor contract address"
}

variable "maker_mock_api_key" {
  type        = string
  description = "Heroku API key"
  default     = "test"
}

variable "maker_mock_is_disabled" {
  type        = string
  description = "Whether to disable the maker mock creating new orders"
  default     = "true"
}

variable "maker_mock_private_key" {
  type        = string
  description = "Private key for the maker mock"
  default     = "test"
}

variable "dd_api_key" {
  type        = string
  description = "Datadog API key"
}
