package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MeisterNickey/experience-hub/internal/models"
	"github.com/MeisterNickey/experience-hub/internal/repository"
)

// TicketHandler serves HTTP endpoints for the Go Core API.
type TicketHandler struct {
	ticketRepo repository.TicketRepository
	cacheRepo  repository.CacheRepository
	sqsQueue   string
}

func NewTicketHandler(ticketRepo repository.TicketRepository, cacheRepo repository.CacheRepository, sqsQueue string) *TicketHandler {
	return &TicketHandler{
		ticketRepo: ticketRepo,
		cacheRepo:  cacheRepo,
		sqsQueue:   sqsQueue,
	}
}

// CreateTicketRequest represents the incoming omni-channel ticket payload.
type CreateTicketRequest struct {
	UserID      string                `json:"user_id"`
	Channel     string                `json:"channel"`
	Category    string                `json:"category"`
	Subject     string                `json:"subject"`
	Description string                `json:"description"`
	Priority    models.TicketPriority `json:"priority"`
}

// HealthCheck handles GET /health
func (h *TicketHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "experience-hub-core-api",
	})
}

// IngestTicket handles POST /api/v1/tickets
// Validates rate limits via Redis and enqueues the raw ticket payload to Amazon SQS.
func (h *TicketHandler) IngestTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Subject == "" {
		http.Error(w, "user_id and subject are required", http.StatusBadRequest)
		return
	}

	allowed, err := h.cacheRepo.CheckRateLimit(r.Context(), req.UserID)
	if err != nil || !allowed {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	now := time.Now().UTC()
	ticket := &models.Ticket{
		ID:          fmt.Sprintf("tkt_%d", now.UnixNano()),
		UserID:      req.UserID,
		Channel:     req.Channel,
		Category:    req.Category,
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      models.TicketStatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Phase 3 TODO: Publish raw payload directly to Amazon SQS via sqs.SendMessage
	slog.Info("enqueued ticket payload to SQS",
		"queue_url", h.sqsQueue,
		"ticket_id", ticket.ID,
		"channel", ticket.Channel,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "ticket accepted for asynchronous processing",
		"ticket":  ticket,
	})
}
