variable "vpc_cidr"        { type = string }
variable "azs"             { type = list(string) }
variable "public_subnets"  { type = list(string) }
# variable "private_subnets" { type = list(string) }
variable "ecr_repos"       { type = list(string) }  # ["backend", "frontend"] など

# 追加: backend イメージタグ（CIでコミットSHAを渡す）
variable "backend_image_tag" {
  type = string
  default = "latest"
}

# DB 等の公開値
variable "db_host" {
  type = string
  default = ""
}
variable "db_user" {
  type = string
  default = ""
}
variable "db_name" {
  type = string
  default = ""
}

# Secrets Manager の名前または ARN をマップで渡す（キーはコンテナ側で使う環境変数名）
# 例: { DB_PASSWORD = "readoku/dev/db_password", JWT_SECRET = "readoku/dev/jwt", GEMINI_API_KEY = "readoku/dev/gemini" }
variable "secrets" {
  type = map(string)
  default = {}
}
