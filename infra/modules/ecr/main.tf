resource "aws_ecr_repository" "this" {
  for_each             = toset(var.repositories)
  name                 = "${var.name_prefix}-${each.key}"
  image_tag_mutability = "IMUTABLE"
  image_scanning_configuration { scan_on_push = true }
  tags = merge(var.tags, { Name = "${var.name_prefix}-${each.key}" })
}

resource "aws_ecr_lifecycle_policy" "this" {
  for_each   = aws_ecr_repository.this
  repository = each.value.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Expire untagged images older than 30 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 30
        }
        action = { type = "expire" }
      },
      {
        rulePriority = 2
        description = "Keep last 10 tagged images"
        selection = {
          tagStatus   = "tagged"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = { type = "expire" }
      }
    ]
  })
}
