package models

import "time"

// TicketStatus represents the lifecycle state of a support ticket.
type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "OPEN"
	TicketStatusAssigned   TicketStatus = "ASSIGNED"
	TicketStatusInProgress TicketStatus = "IN_PROGRESS"
	TicketStatusResolved   TicketStatus = "RESOLVED"
	TicketStatusClosed     TicketStatus = "CLOSED"
)

// TicketPriority represents urgency level for routing.
type TicketPriority string

const (
	PriorityLow      TicketPriority = "LOW"
	PriorityMedium   TicketPriority = "MEDIUM"
	PriorityHigh     TicketPriority = "HIGH"
	PriorityCritical TicketPriority = "CRITICAL"
)

// Ticket stores structured relational state in Amazon RDS (PostgreSQL).
type Ticket struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	AgentID     *string        `json:"agent_id,omitempty"`
	Channel     string         `json:"channel"` // e.g., "email", "chat", "web", "whatsapp"
	Category    string         `json:"category"`
	Subject     string         `json:"subject"`
	Description string         `json:"description"`
	Priority    TicketPriority `json:"priority"`
	Status      TicketStatus   `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// User represents a customer profile stored in Amazon RDS (PostgreSQL).
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Tier      string    `json:"tier"` // e.g., "free", "pro", "enterprise"
	CreatedAt time.Time `json:"created_at"`
}

// Agent represents a live support agent whose availability is cached in Redis.
type Agent struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	Skills      []string  `json:"skills"`
	IsAvailable bool      `json:"is_available"`
	ActiveCount int       `json:"active_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AuditLog represents high-velocity unstructured events stored in Amazon DynamoDB.
// Partition Key: TicketID, Sort Key: Timestamp.
type AuditLog struct {
	TicketID  string                 `json:"TicketID"`
	Timestamp string                 `json:"Timestamp"`
	EventType string                 `json:"event_type"` // e.g., "TICKET_INGESTED", "AGENT_ROUTED", "TRANSCRIPT_APPENDED"
	ActorID   string                 `json:"actor_id"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
