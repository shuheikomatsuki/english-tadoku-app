module "networking" {
  source      = "./modules/networking"
  name_prefix = local.name_prefix
  tags        = local.tags
}

module "frontend" {
  source      = "./modules/frontend"
  name_prefix = local.name_prefix
  tags        = local.tags
}

module "backend" {
  source      = "./modules/backend"
  name_prefix = local.name_prefix
  tags        = local.tags
  project_name           = var.project_name
  environment            = var.environment
  enable_backend         = var.enable_backend
  lambda_package_path    = var.backend_lambda_package_path
  frontend_url           = var.backend_frontend_url
  allowed_origins        = var.backend_allowed_origins
  daily_generation_limit = var.backend_daily_generation_limit
  parameter_prefix_override = null
  log_retention_in_days     = 14
  lambda_memory_size        = 256
  lambda_timeout_seconds    = 30
  lambda_architectures      = ["arm64"]
}
