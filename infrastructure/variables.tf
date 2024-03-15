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

variable "verify_signature" {
  type        = string
  description = "Whether to verify the signature sent with the order (not compulsory as the signature is verified on-chain)"
  default     = "false"
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
