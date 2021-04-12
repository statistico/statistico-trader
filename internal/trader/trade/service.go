package trade

import (
	"context"
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
)

type Writer interface {
	Insert(t *Trade) error
}

type Reader interface {
	Get(q *ReaderQuery) ([]*Trade, error)
	Exists(market, runner string, eventID uint64, strategyID uuid.UUID) (bool, error)
}

type ReaderQuery struct {
	StrategyID   uuid.UUID
	Result       []string
}

type Placer interface {
	// PlaceTrade receives an exchange.Client struct to place a Trade record with an external exchange and returns
	// the resulting Trade struct.
	PlaceTrade(ctx context.Context, c exchange.Client, t *Ticket, s *strategy.Strategy) (*Trade, error)
}

type Manager interface {
	Manage(ctx context.Context, t *Ticket, s *strategy.Strategy) error
}

