package repository

import (
	"context"
	"sync"

	"github.com/MeisterNickey/experience-hub/internal/models"
)

// TicketRepository defines relational persistence operations for Amazon RDS (PostgreSQL).
type TicketRepository interface {
	CreateTicket(ctx context.Context, ticket *models.Ticket) error
	GetTicketByID(ctx context.Context, id string) (*models.Ticket, error)
	UpdateTicketStatus(ctx context.Context, id string, status models.TicketStatus, agentID *string) error
}

// PostgresRepository implements TicketRepository backed by PostgreSQL (RDS).
// Phase 3 TODO: Replace in-memory fallback with pgxpool / GORM connection.
type PostgresRepository struct {
	dsn     string
	mu      sync.RWMutex
	tickets map[string]*models.Ticket
}

func NewPostgresRepository(dsn string) *PostgresRepository {
	return &PostgresRepository{
		dsn:     dsn,
		tickets: make(map[string]*models.Ticket),
	}
}

func (r *PostgresRepository) CreateTicket(ctx context.Context, ticket *models.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tickets[ticket.ID] = ticket
	return nil
}

func (r *PostgresRepository) GetTicketByID(ctx context.Context, id string) (*models.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tickets[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *PostgresRepository) UpdateTicketStatus(ctx context.Context, id string, status models.TicketStatus, agentID *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.tickets[id]; ok {
		t.Status = status
		t.AgentID = agentID
	}
	return nil
}
