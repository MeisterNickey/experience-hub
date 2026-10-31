package repository

import (
	"context"
	"sync"
)

// CacheRepository handles rate limiting and live agent availability in ElastiCache Redis.
type CacheRepository interface {
	CheckRateLimit(ctx context.Context, clientID string) (bool, error)
	GetAvailableAgent(ctx context.Context, category string) (string, error)
	SetAgentAvailability(ctx context.Context, agentID string, available bool) error
}

// RedisRepository implements CacheRepository backed by Amazon ElastiCache (Redis).
// Phase 3 TODO: Wire go-redis client for distributed rate limiting & agent routing sets.
type RedisRepository struct {
	addr   string
	mu     sync.RWMutex
	agents map[string]bool
}

func NewRedisRepository(addr, password string) *RedisRepository {
	return &RedisRepository{
		addr: addr,
		agents: map[string]bool{
			"agent-101": true,
			"agent-102": true,
		},
	}
}

func (r *RedisRepository) CheckRateLimit(ctx context.Context, clientID string) (bool, error) {
	return true, nil
}

func (r *RedisRepository) GetAvailableAgent(ctx context.Context, category string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for id, available := range r.agents {
		if available {
			return id, nil
		}
	}
	return "", nil
}

func (r *RedisRepository) SetAgentAvailability(ctx context.Context, agentID string, available bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[agentID] = available
	return nil
}
