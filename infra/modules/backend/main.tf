locals {
  ssm_parameter_prefix = coalesce(
    var.parameter_prefix_override,
    "/${var.project_name}/${var.environment}/",
  )

  dev_cors_allow_origins = var.environment == "dev" ? [
    "http://localhost:5173",
    "http://127.0.0.1:5173",
  ] : []
  
  cors_allow_origins = distinct(
    compact(
      concat(
        var.allowed_origins,
        [
          var.frontend_url,
        ],
        local.dev_cors_allow_origins,
      )
    )
  )

  cors_allow_headers = [
    "Content-Type",
    "Authorization",
    "Accept",
    "Origin",
    "X-Requested-With",
  ]

  lambda_env = {
    FRONTEND_URL            = var.frontend_url
    DAILY_GENERATION_LIMIT  = tostring(var.daily_generation_limit)
    GEMINI_API_KEY_PARAM    = "${local.ssm_parameter_prefix}gemini_api_key"
    PARAMETER_PREFIX        = local.ssm_parameter_prefix
    ENVIRONMENT             = var.environment
    DB_HOST_PARAM           = "${local.ssm_parameter_prefix}db_host"
    DB_USER_PARAM           = "${local.ssm_parameter_prefix}db_user"
    DB_PASSWORD_PARAM       = "${local.ssm_parameter_prefix}db_password"
    DB_NAME_PARAM           = "${local.ssm_parameter_prefix}db_name"
    JWT_SECRET_PARAM        = "${local.ssm_parameter_prefix}jwt_secret"
  }
}

data "aws_caller_identity" "current" {
  count = var.enable_backend ? 1 : 0
}

data "aws_region" "current" {
  count = var.enable_backend ? 1 : 0
}

resource "aws_iam_role" "lambda" {
  count = var.enable_backend ? 1 : 0

  name               = "${var.name_prefix}-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role[0].json
  tags               = var.tags
}

data "aws_iam_policy_document" "lambda_assume_role" {
  count = var.enable_backend ? 1 : 0

  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  count      = var.enable_backend ? 1 : 0
  role       = aws_iam_role.lambda[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "ssm_read" {
  count = var.enable_backend ? 1 : 0

  statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
    ]

    resources = [
      "arn:aws:ssm:${data.aws_region.current[0].name}:${data.aws_caller_identity.current[0].account_id}:parameter${local.ssm_parameter_prefix}*",
    ]
  }
}

resource "aws_iam_policy" "ssm_read" {
  count = var.enable_backend ? 1 : 0

  name        = "${var.name_prefix}-ssm-read"
  description = "Allow Lambda to read parameters under ${local.ssm_parameter_prefix}"
  policy      = data.aws_iam_policy_document.ssm_read[0].json
}

resource "aws_iam_role_policy_attachment" "lambda_ssm_read" {
  count      = var.enable_backend ? 1 : 0
  role       = aws_iam_role.lambda[0].name
  policy_arn = aws_iam_policy.ssm_read[0].arn
}

resource "aws_cloudwatch_log_group" "lambda" {
  count = var.enable_backend ? 1 : 0

  name              = "/aws/lambda/${var.name_prefix}-api"
  retention_in_days = var.log_retention_in_days
  tags              = var.tags
}

resource "aws_lambda_function" "api" {
  count = var.enable_backend ? 1 : 0

  function_name = "${var.name_prefix}-api"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  role          = aws_iam_role.lambda[0].arn
  filename      = var.lambda_package_path
  source_code_hash = filebase64sha256(var.lambda_package_path)

  architectures = var.lambda_architectures
  memory_size   = var.lambda_memory_size
  timeout       = var.lambda_timeout_seconds

  environment {
    variables = local.lambda_env
  }

  tags = var.tags
}

resource "aws_apigatewayv2_api" "api" {
  count = var.enable_backend ? 1 : 0

  name          = "${var.name_prefix}-http-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_credentials = true
    allow_origins     = local.cors_allow_origins
    allow_methods     = ["GET", "POST", "PATCH", "DELETE", "OPTIONS"]
    allow_headers     = local.cors_allow_headers
  }
}

resource "aws_apigatewayv2_integration" "lambda" {
  count = var.enable_backend ? 1 : 0

  api_id                 = aws_apigatewayv2_api.api[0].id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.api[0].invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "default" {
  count = var.enable_backend ? 1 : 0

  api_id    = aws_apigatewayv2_api.api[0].id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.lambda[0].id}"
}

resource "aws_apigatewayv2_stage" "default" {
  count = var.enable_backend ? 1 : 0

  api_id      = aws_apigatewayv2_api.api[0].id
  name        = "$default"
  auto_deploy = true
  tags        = var.tags
}

resource "aws_lambda_permission" "apigw" {
  count = var.enable_backend ? 1 : 0

  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api[0].function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api[0].execution_arn}/*/*"
}
