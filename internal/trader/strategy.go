package trader

import "github.com/google/uuid"

type StrategyWriter interface {
	Insert(s *Strategy) error
}

type StrategyReader interface {
	Get(q *StrategyReaderQuery) ([]*Strategy, error)
}

type StrategyReaderQuery struct {
	UserID   *uuid.UUID
	Visibility  string
}
