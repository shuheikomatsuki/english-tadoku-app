output "bucket_id" {
  description = "S3 bucket name for the frontend."
  value       = aws_s3_bucket.site.id
}

output "bucket_arn" {
  description = "S3 bucket ARN for the frontend."
  value       = aws_s3_bucket.site.arn
}

output "cloudfront_domain_name" {
  description = "CloudFront domain name for the distribution."
  value       = aws_cloudfront_distribution.site.domain_name
}

output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID."
  value       = aws_cloudfront_distribution.site.id
}
