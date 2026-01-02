variable "name_prefix" {
  description = "Prefix used for naming resources."
  type        = string
}

variable "tags" {
  description = "Common tags to apply to resources."
  type        = map(string)
  default     = {}
}

variable "project_name" {
  description = "Project name used for naming and parameter paths."
  type        = string
}

variable "environment" {
  description = "Environment name (e.g. dev, prod)."
  type        = string
}

variable "enable_backend" {
  description = "Whether to create backend resources."
  type        = bool
  default     = false
}

variable "lambda_package_path" {
  description = "Path to the zipped Lambda package."
  type        = string
  default     = ""
  validation {
    condition = var.enable_backend == false || length(trimspace(var.lambda_package_path)) > 0
    error_message = "When enable_backend is true, lambda_package_path must be a non-empty path to the zipped Lambda package."
  }
}

variable "lambda_memory_size" {
  description = "Lambda memory size in MB."
  type        = number
  default     = 256
}

variable "lambda_timeout_seconds" {
  description = "Lambda timeout in seconds."
  type        = number
  default     = 30
}

variable "lambda_architectures" {
  description = "Lambda architectures (e.g. [\"arm64\"], [\"x86_64\"])."
  type        = list(string)
  default     = ["arm64"]
}

variable "frontend_url" {
  description = "Primary frontend origin for CORS."
  type        = string
  default     = "http://localhost:5173"
}

variable "allowed_origins" {
  description = "Additional allowed origins for CORS."
  type        = list(string)
  default     = []
}

variable "daily_generation_limit" {
  description = "Daily generation limit passed to the backend."
  type        = number
  default     = 10
}

variable "parameter_prefix_override" {
  description = "Optional override for SSM parameter prefix (default: /<project>/<env>/)."
  type        = string
  default     = null
}

variable "log_retention_in_days" {
  description = "CloudWatch Logs retention in days."
  type        = number
  default     = 14
}
