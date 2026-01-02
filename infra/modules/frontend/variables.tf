variable "name_prefix" {
  description = "Prefix used for naming resources."
  type        = string
}

variable "tags" {
  description = "Common tags to apply to resources."
  type        = map(string)
  default     = {}
}
