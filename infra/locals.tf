locals {
  project     = "readoku"
  environment = "dev"
  name_prefix = "${local.project}-${local.environment}"
  tags = {
    Project   = local.project
    Env       = local.environment
    ManagedBy = "terraform"
  }
}
