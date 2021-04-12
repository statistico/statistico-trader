package strategy

import (
	"context"
	"github.com/google/uuid"
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

type FinderQuery struct {
	MarketName    string    `json:"marketName"`
	RunnerName    string    `json:"runnerName"`
	EventID       uint64    `json:"eventId"`
	CompetitionID uint64    `json:"competitionId"`
	Price         float32   `json:"price"`
	Side          string    `json:"side"`
	Status        string    `json:"status"`
}

type Finder interface {
	FindMatchingStrategies(ctx context.Context, q *FinderQuery) <-chan *Strategy
}
