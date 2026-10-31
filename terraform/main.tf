# ------------------------------------------------------------------------------
# 1. Messaging: Amazon SQS Queue + Dead Letter Queue (DLQ)
# ------------------------------------------------------------------------------
resource "aws_sqs_queue" "tickets_dlq" {
  name                      = "${var.project_name}-${var.environment}-tickets-dlq"
  message_retention_seconds = 1209600 # 14 days
}

resource "aws_sqs_queue" "tickets_queue" {
  name                       = "${var.project_name}-${var.environment}-tickets-queue"
  visibility_timeout_seconds = 60
  message_retention_seconds  = 345600 # 4 days
  receive_wait_time_seconds  = 20     # Enable long polling

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.tickets_dlq.arn
    maxReceiveCount     = 5
  })
}

# ------------------------------------------------------------------------------
# 2. Unstructured Event Store: Amazon DynamoDB (Partition Key: TicketID)
# ------------------------------------------------------------------------------
resource "aws_dynamodb_table" "audit_logs" {
  name         = "${var.project_name}-${var.environment}-audit-logs"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "TicketID"
  range_key    = "Timestamp"

  attribute {
    name = "TicketID"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }
}

# ------------------------------------------------------------------------------
# 3. Relational State Store: Amazon RDS (PostgreSQL)
# ------------------------------------------------------------------------------
resource "aws_db_instance" "postgres" {
  identifier             = "${var.project_name}-${var.environment}-postgres"
  engine                 = "postgres"
  engine_version         = "16.3"
  instance_class         = "db.t4g.micro"
  allocated_storage      = 20
  max_allocated_storage  = 100
  storage_encrypted      = true
  db_name                = "experience_hub"
  username               = var.db_username
  password               = var.db_password
  skip_final_snapshot    = true
  publicly_accessible    = false
  deletion_protection    = false
}

# ------------------------------------------------------------------------------
# 4. Cache & Rate Limiting: Amazon ElastiCache (Redis)
# ------------------------------------------------------------------------------
resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "${var.project_name}-${var.environment}-redis"
  engine               = "redis"
  node_type            = "cache.t4g.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  engine_version       = "7.1"
  port                 = 6379
}
