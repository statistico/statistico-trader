package trader

import "github.com/google/uuid"

type TradeWriter interface {
	Insert(t *Trade) error
}

type TradeReader interface {
	Get(q *TradeReaderQuery) ([]*Trade, error)
}

type TradeReaderQuery struct {
	StrategyID   uuid.UUID
	Result       []string
}
