package trade

import "github.com/google/uuid"

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
