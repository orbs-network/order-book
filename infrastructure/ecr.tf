resource "aws_ecr_repository" "my_repo" {
  name                 = "order-book-repo-${var.environment_name}"
  image_tag_mutability = "IMMUTABLE"
  tags                 = local.tags
}


resource "aws_ecr_lifecycle_policy" "this" {
  repository = aws_ecr_repository.my_repo.name

  policy = <<EOF
{
    "rules": [
        {
            "rulePriority": 1,
            "description": "Keep last 100 images",
            "selection": {
                "tagStatus": "any",
                "countType": "imageCountMoreThan",
                "countNumber": 100
            },
            "action": {
                "type": "expire"
            }
        }
    ]
}
EOF
}
