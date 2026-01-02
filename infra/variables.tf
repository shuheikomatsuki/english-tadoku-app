variable "project_name" {
  description = "Short name of the project (used for tagging and naming)."
  type        = string
}

variable "environment" {
  description = "Deployment environment identifier, e.g. dev or prod."
  type        = string
}

variable "aws_region" {
  description = "AWS region to deploy into."
  type        = string
}

variable "aws_profile" {
  description = "Optional AWS CLI profile name (null to use default credentials chain)."
  type        = string
  default     = null
}

variable "enable_backend" {
  description = "Whether to create backend (Lambda + API Gateway) resources."
  type        = bool
  default     = false
}

variable "backend_lambda_package_path" {
  description = "Path to the zipped Lambda package for the backend. Used when enable_backend is true."
  type        = string
  default     = ""
}

variable "backend_frontend_url" {
  description = "Frontend origin used for CORS in API Gateway."
  type        = string
  default     = "http://localhost:5173"
}

variable "backend_allowed_origins" {
  description = "Additional allowed origins for API Gateway CORS (frontend_url and localhost are included automatically)."
  type        = list(string)
  default     = []
}

variable "backend_daily_generation_limit" {
  description = "Daily generation limit to pass to the backend service."
  type        = number
  default     = 10
}

variable "backend_log_retention_in_days" {
  description = "CloudWatch log retention for backend Lambda."
  type        = number
  default     = 14
}

variable "backend_lambda_memory_size" {
  description = "Memory size (MB) for backend Lambda."
  type        = number
  default     = 256
}

variable "backend_lambda_timeout_seconds" {
  description = "Timeout (seconds) for backend Lambda."
  type        = number
  default     = 30
}

variable "backend_lambda_architectures" {
  description = "Architectures for backend Lambda."
  type        = list(string)
  default     = ["arm64"]
}
