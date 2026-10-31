package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/MeisterNickey/experience-hub/internal/models"
	"github.com/MeisterNickey/experience-hub/internal/repository"
)

// Processor polls Amazon SQS using concurrent Goroutines and coordinates
// writes across Amazon RDS (PostgreSQL), Amazon DynamoDB, and ElastiCache (Redis).
type Processor struct {
	queueURL    string
	concurrency int
	ticketRepo  repository.TicketRepository
	auditRepo   repository.AuditLogRepository
	cacheRepo   repository.CacheRepository
}

func NewProcessor(
	queueURL string,
	concurrency int,
	ticketRepo repository.TicketRepository,
	auditRepo repository.AuditLogRepository,
	cacheRepo repository.CacheRepository,
) *Processor {
	return &Processor{
		queueURL:    queueURL,
		concurrency: concurrency,
		ticketRepo:  ticketRepo,
		auditRepo:   auditRepo,
		cacheRepo:   cacheRepo,
	}
}

// Start launches the worker pool of Goroutines polling SQS until ctx is cancelled.
func (p *Processor) Start(ctx context.Context) {
	var wg sync.WaitGroup
	slog.Info("starting SQS ticket processor pool",
		"queue_url", p.queueURL,
		"concurrency", p.concurrency,
	)

	for i := 1; i <= p.concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			p.pollLoop(ctx, workerID)
		}(i)
	}

	wg.Wait()
	slog.Info("all SQS worker goroutines stopped cleanly")
}

func (p *Processor) pollLoop(ctx context.Context, workerID int) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping worker goroutine", "worker_id", workerID)
			return
		case <-ticker.C:
			// Phase 3 TODO: Replace ticker with SQS ReceiveMessage long-polling (WaitTimeSeconds=20)
			slog.Debug("polling SQS queue for ticket messages", "worker_id", workerID, "queue", p.queueURL)
		}
	}
}

// ProcessTicket handles a single dequeued ticket:
// 1. Looks up an available agent in Redis
// 2. Writes relational state to PostgreSQL (RDS)
// 3. Appends unstructured event log to DynamoDB
func (p *Processor) ProcessTicket(ctx context.Context, ticket *models.Ticket) error {
	agentID, err := p.cacheRepo.GetAvailableAgent(ctx, ticket.Category)
	if err == nil && agentID != "" {
		ticket.AgentID = &agentID
		ticket.Status = models.TicketStatusAssigned
	}

	if err := p.ticketRepo.CreateTicket(ctx, ticket); err != nil {
		return err
	}

	audit := &models.AuditLog{
		TicketID:  ticket.ID,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		EventType: "TICKET_PROCESSED",
		ActorID:   "system-worker",
		Metadata: map[string]interface{}{
			"channel":  ticket.Channel,
			"category": ticket.Category,
			"status":   ticket.Status,
			"agent_id": ticket.AgentID,
		},
	}

	return p.auditRepo.PutAuditLog(ctx, audit)
}
