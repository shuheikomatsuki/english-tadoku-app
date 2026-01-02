output "lambda_function_name" {
  description = "Backend Lambda function name."
  value       = var.enable_backend ? aws_lambda_function.api[0].function_name : null
}

output "lambda_role_arn" {
  description = "IAM role ARN for the backend Lambda."
  value       = var.enable_backend ? aws_iam_role.lambda[0].arn : null
}

output "api_endpoint" {
  description = "HTTP API endpoint URL."
  value       = var.enable_backend ? aws_apigatewayv2_api.api[0].api_endpoint : null
}
