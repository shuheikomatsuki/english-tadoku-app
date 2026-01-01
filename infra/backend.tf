module "vpc" {
  source          = "./modules/vpc"
  name            = "${local.name_prefix}-vpc"
  cidr            = var.vpc_cidr
  azs             = var.azs
  public_subnets  = var.public_subnets
#   private_subnets = var.private_subnets
  tags            = local.tags
}

module "ecr" {
  source       = "./modules/ecr"
  name_prefix  = local.name_prefix
  repositories = var.ecr_repos
  tags         = local.tags
}
