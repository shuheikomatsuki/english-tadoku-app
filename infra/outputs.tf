output "frontend_bucket_id" {
  description = "S3 bucket name for the frontend."
  value       = module.frontend.bucket_id
}

output "frontend_bucket_arn" {
  description = "S3 bucket ARN for the frontend."
  value       = module.frontend.bucket_arn
}

output "frontend_cloudfront_domain_name" {
  description = "CloudFront domain name for the frontend."
  value       = module.frontend.cloudfront_domain_name
}

output "frontend_cloudfront_distribution_id" {
  description = "CloudFront distribution ID for the frontend."
  value       = module.frontend.cloudfront_distribution_id
}
