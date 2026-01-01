# infra/

Terraform による AWS インフラ定義（dev環境）

## 前提
- Terraform >= 1.5.0
- AWS CLI 設定済み
- region: ap-northeast-1

## ディレクトリ構成
```
infra/
├─ README.md
├─ versions.tf
├─ provider.tf
├─ backend.tf
├─ variables.tf
├─ outputs.tf
├─ locals.tf
├─ modules/
│  ├─ vpc/
│  │  ├─ main.tf
│  │  ├─ variables.tf
│  │  └─ outputs.tf
│  ├─ ecr/
│  │  ├─ main.tf
│  │  └─ outputs.tf
│  ├─ ecs/
│  │  ├─ cluster.tf
│  │  ├─ task_definition.tf
│  │  ├─ service.tf
│  │  └─ iam.tf
│  ├─ alb/
│  │  ├─ main.tf
│  │  └─ outputs.tf
│  ├─ rds/
│  │  ├─ main.tf
│  │  ├─ subnet_group.tf
│  │  └─ outputs.tf
│  ├─ secrets/
│  │  └─ main.tf
│  ├─ s3_cloudfront/
│  │  ├─ s3.tf
│  │  ├─ cloudfront.tf
│  │  └─ outputs.tf
│  ├─ route53/
│  │  └─ main.tf
│  └─ iam/
│     └─ main.tf
├─ envs/
│  └─ dev/
│     ├─ main.tf
│     ├─ terraform.tfvars
│     └─ provider_override.tf   # 区別用（optional）
├─ scripts/
│  ├─ build_and_push.sh
│  └─ deploy_helper.sh
└─ .gitignore
```
