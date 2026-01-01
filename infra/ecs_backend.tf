resource "aws_ecs_cluster" "backend" {
  name = "${local.name_prefix}-cluster"
  tags = local.tags
}

data "aws_iam_policy_document" "ecs_task_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_task_execution" {
  name               = "${local.name_prefix}-ecs-exec"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = local.tags
}

resource "aws_iam_role_policy_attachment" "ecs_task_exec_policy" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# 追加: アプリ用のタスクロール（最小権限を付ける。必要なポリシーのみ添付すること）
resource "aws_iam_role" "ecs_task_role" {
  name               = "${local.name_prefix}-ecs-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = local.tags
}

# 例: Secrets Manager を読み取る権限を与える（必要なら専用の最小ポリシーを作ることを推奨）
resource "aws_iam_role_policy_attachment" "task_secrets_read" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = "arn:aws:iam::aws:policy/SecretsManagerReadWrite" # 開発用。運用では最小権限ポリシーへ差し替え
}

resource "aws_security_group" "backend" {
  name        = "${local.name_prefix}-backend-sg"
  description = "Allow HTTP to Fargate task"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

# Secrets Manager データを for_each で参照
data "aws_secretsmanager_secret" "secrets" {
  for_each = var.secrets
  name     = each.value
}
data "aws_secretsmanager_secret_version" "secrets" {
  for_each  = var.secrets
  secret_id = data.aws_secretsmanager_secret.secrets[each.key].id
}

resource "aws_ecs_task_definition" "backend" {
  family                   = "${local.name_prefix}-backend"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([{
    name         = "backend"
    image        = "${module.ecr.repository_urls["backend"]}:${var.backend_image_tag}"
    essential    = true
    portMappings = [{ containerPort = 8080, hostPort = 80, protocol = "tcp" }]
    # 公開値は environment で渡す
    environment = [
      { name = "DB_HOST", value = var.db_host },
      { name = "DB_USER", value = var.db_user },
      { name = "DB_NAME", value = var.db_name }
    ]

    # 機密は Secrets Manager を参照して secrets で渡す
    secrets = [
      for k in keys(var.secrets) : {
        name      = k
        valueFrom = data.aws_secretsmanager_secret_version.secrets[k].arn
      }
    ]
  }])

  tags = local.tags
}

# とりあえず今は最小構成で動くか検証するのが優先

# # 追加: Cloud Map namesapce / service（サービスディスカバリ）
# resource "aws_service_discovery_private_dns_namespace" "ns" {
#   name = "${local.name_prefix}.local"
#   vpc  = module.vpc.vpc_id
# }

# resource "aws_service_discovery_service" "backend" {
#   name        = "${local.name_prefix}-backend-sd"
#   namespace_id = aws_service_discovery_private_dns_namespace.ns.id

#   dns_config {
#     namespace_id = aws_service_discovery_private_dns_namespace.ns.id
#     dns_records {
#       ttl = 10
#       type = "A"
#     }
#     routing_policy = "MULTIVALUE"
#   }

#   health_check_custom_config { failure_threshold = 1 }
# }

resource "aws_ecs_service" "backend" {
  name            = "${local.name_prefix}-backend"
  cluster         = aws_ecs_cluster.backend.id
  task_definition = aws_ecs_task_definition.backend.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    assign_public_ip = true
    subnets          = module.vpc.public_subnet_ids
    security_groups  = [aws_security_group.backend.id]
  }

#   service_registries {
#     registry_arn = aws_service_discovery_service.backend.arn
#   }

  tags = local.tags
}
