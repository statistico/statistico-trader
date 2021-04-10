package strategy

import (
	"context"
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/market"
)

type Writer interface {
	Insert(s *Strategy) error
}

type Reader interface {
	Get(q *ReaderQuery) ([]*Strategy, error)
}

type ReaderQuery struct {
	UserID     *uuid.UUID
	Market     *string
	Runner     *string
	Price      *float32
	CompetitionID *uint64
	Side       *string
	Status     *string
	Visibility *string
	OrderBy    *string
}

type Finder interface {
	AddMatchingStrategies(ctx context.Context, m *market.Runner, ch chan<- *Strategy)
}
