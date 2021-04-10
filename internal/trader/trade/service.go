package trade

import (
	"context"
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/market"
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

type Manager interface {
	PlaceTrade(ctx context.Context, c exchange.Client, r *market.Runner, s strategy.Strategy) (*Trade, error)
}
