variable "vpc_cidr"        { type = string }
variable "azs"             { type = list(string) }
variable "public_subnets"  { type = list(string) }
# variable "private_subnets" { type = list(string) }
variable "ecr_repos"       { type = list(string) }  # ["backend", "frontend"] など
