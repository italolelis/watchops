# ============================
# Kinesis KMS
# ============================
resource "aws_kms_key" "this" {
  deletion_window_in_days = var.kms_deletion_window_in_days
  enable_key_rotation     = true
}

resource "aws_kms_alias" "this" {
  name          = format("alias/%s-kinesis-enc-key", var.name)
  target_key_id = aws_kms_key.this.key_id
}

# ============================
# Kinesis Stream
# Creates a kinesis github stream that will hold all events related to github
# ============================
resource "aws_kinesis_stream" "this" {
  name                      = var.name
  shard_count               = var.shard_count
  retention_period          = var.retention_period
  shard_level_metrics       = var.shard_level_metrics
  enforce_consumer_deletion = var.enforce_consumer_deletion
  tags                      = var.tags

  encryption_type = "KMS"
  kms_key_id      = aws_kms_key.this.arn

  lifecycle {
    ignore_changes = [shard_count]
  }
}

# ============================
# Kinesis policy read only
# ============================
resource "aws_iam_policy" "read-only" {
  count = var.create_policy_read_only == true ? 1 : 0

  name        = format("kinesis-stream-%s-read-only", var.name)
  path        = "/"
  description = format("Kinesis %s read-only policy", var.name)
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = concat([
      {
        Effect = "Allow"
        Action = [
          "kinesis:DescribeLimits",
          "kinesis:DescribeStream",
          "kinesis:GetRecords",
          "kinesis:GetShardIterator",
          "kinesis:SubscribeToShard",
          "kinesis:ListShards"
        ]
        Resource = [
          aws_kinesis_stream.this.arn
        ]
      }
    ])
  })
}

# ============================
# Kinesis policy write only
# ============================
resource "aws_iam_policy" "write-only" {
  count = var.create_policy_write_only == true ? 1 : 0

  name        = format("kinesis-stream-%s-write-only", var.name)
  path        = "/"
  description = format("Kinesis %s write-only policy", var.name)
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = concat([
      {
        Effect = "Allow"
        Action = [
          "kinesis:DescribeStream",
          "kinesis:PutRecord",
          "kinesis:PutRecords",
        ]
        Resource = [
          aws_kinesis_stream.this.arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "kms:Encrypt",
          "kms:DescribeKey",
          "kms:GenerateDataKey",
          "kms:Decrypt"
        ]
        Resource = [
          aws_kms_key.this.arn
        ]
      }
    ])
  })
}
