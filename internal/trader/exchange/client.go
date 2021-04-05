package exchange

import "time"

type Client interface {
	Balance() (float32, error)
	PlaceTrade(t *TradeTicket) (*Trade, error)
}

type TradeTicket struct {
	MarketID  string
	RunnerID  uint64
	Price     float32
	Stake     float32
	Side      string
	OrderType string
	PersistenceType string
}

type Trade struct {
	Reference string
	Status  string
	Timestamp time.Time
}
