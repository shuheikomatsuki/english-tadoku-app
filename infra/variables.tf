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
