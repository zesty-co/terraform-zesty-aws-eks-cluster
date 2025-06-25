resource "random_uuid" "zesty_external_id" {}

data "aws_caller_identity" "current" {}

resource "aws_iam_role" "zesty_iam_role" {
  name                 = var.role_name
  max_session_duration = var.max_session_duration

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Principal = {
          AWS = var.trusted_principal
        }
        Condition = {
          StringEquals = {
            "sts:ExternalId" = random_uuid.zesty_external_id.result
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "zesty_policy" {
  name = "Zesty-Policy-${data.aws_caller_identity.current.account_id}"
  role = aws_iam_role.zesty_iam_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "EC2Access"
        Effect = "Allow"
        Action = [
          "ec2:List*",
          "ec2:Describe*",
          "elasticloadbalancing:Describe*",
          "autoscaling:Describe*"
        ]
        Resource = ["*"]
      },
      {
        Sid    = "OrganizationsAccess"
        Effect = "Allow"
        Action = [
          "organizations:List*",
          "organizations:Describe*"
        ]
        Resource = ["*"]
      },
      {
        Sid    = "ServiceQuotasAccess"
        Effect = "Allow"
        Action = [
          "servicequotas:ListServiceQuotas",
          "servicequotas:GetServiceQuota",
          "servicequotas:GetRequestedServiceQuotaChange"
        ]
        Resource = ["*"]
      },
      {
        Sid    = "MetricsAccess"
        Effect = "Allow"
        Action = [
          "cloudwatch:List*",
          "cloudwatch:Describe*",
          "cloudwatch:GetMetricStatistics"
        ]
        Resource = ["*"]
      },
      {
        Sid    = "SavingsPlansAccess"
        Effect = "Allow"
        Action = [
          "savingsplans:List*",
          "savingsplans:Describe*"
        ]
        Resource = ["*"]
      },
      {
        Sid    = "CostExplorerAccess"
        Effect = "Allow"
        Action = [
          "ce:List*",
          "ce:Describe*",
          "ce:Get*"
        ]
        Resource = ["*"]
      },
      {
        Sid    = "EKSAccess"
        Effect = "Allow"
        Action = [
          "eks:List*",
          "eks:Describe*"
        ]
        Resource = ["*"]
      }
    ]
  })
}

resource "zesty_account" "result" {
  account = {
    id             = data.aws_caller_identity.current.account_id
    cloud_provider = var.cloud_provider
    role_arn       = aws_iam_role.zesty_iam_role.arn
    external_id    = random_uuid.zesty_external_id.result
    products       = var.products
  }
  depends_on = [aws_iam_role_policy.zesty_policy]
}

resource "local_file" "kompass_values" {
  content = [
    for p in zesty_account.result.account.products : p.values
    if p.name == "Kompass" && p.active == true
  ][0]
  filename   = var.values_yaml_filename
  depends_on = [zesty_account.result]
}
