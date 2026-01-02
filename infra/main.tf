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
}
