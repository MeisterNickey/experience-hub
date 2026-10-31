output "sqs_queue_url" {
  description = "URL of the SQS ticket ingestion queue"
  value       = aws_sqs_queue.tickets_queue.id
}

output "sqs_dlq_url" {
  description = "URL of the SQS Dead Letter Queue"
  value       = aws_sqs_queue.tickets_dlq.id
}

output "dynamodb_table_name" {
  description = "Name of the DynamoDB audit logs table"
  value       = aws_dynamodb_table.audit_logs.name
}

output "rds_endpoint" {
  description = "Connection endpoint for the RDS PostgreSQL instance"
  value       = aws_db_instance.postgres.endpoint
}

output "redis_endpoint" {
  description = "Primary endpoint for the ElastiCache Redis cluster"
  value       = aws_elasticache_cluster.redis.cache_nodes[0].address
}
