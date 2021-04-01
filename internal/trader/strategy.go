package trader

import "github.com/google/uuid"

type StrategyWriter interface {
	Insert(s *Strategy) error
}

type StrategyReader interface {
	Get(q *StrategyReaderQuery) ([]*Strategy, error)
}

type StrategyReaderQuery struct {
	UserID     *uuid.UUID
	Market     *string
	Runner     *string
	Price      *float32
	CompetitionID *uint64
	Visibility *string
	OrderBy    *string
}
