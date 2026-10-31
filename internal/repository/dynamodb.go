package repository

import (
	"context"
	"log/slog"

	"github.com/MeisterNickey/experience-hub/internal/models"
)

// AuditLogRepository defines unstructured event logging for Amazon DynamoDB.
type AuditLogRepository interface {
	PutAuditLog(ctx context.Context, entry *models.AuditLog) error
}

// DynamoDBRepository implements AuditLogRepository using a table keyed by TicketID.
// Phase 3 TODO: Wire AWS SDK for Go v2 dynamodb.Client PutItem call.
type DynamoDBRepository struct {
	tableName string
	region    string
}

func NewDynamoDBRepository(region, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		region:    region,
		tableName: tableName,
	}
}

func (r *DynamoDBRepository) PutAuditLog(ctx context.Context, entry *models.AuditLog) error {
	slog.Info("persisting audit log to DynamoDB",
		"table", r.tableName,
		"ticket_id", entry.TicketID,
		"event_type", entry.EventType,
	)
	return nil
}
