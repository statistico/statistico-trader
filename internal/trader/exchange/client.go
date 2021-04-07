package exchange

import (
	"context"
)

type Client interface {
	Account(ctx context.Context) (*Account, error)
	PlaceTrade(ctx context.Context, t *TradeTicket) (*Trade, error)
}

type Account struct {
	Balance       float32
	Exposure      float32
	ExposureLimit float32
}

type TradeTicket struct {
	MarketID        string
	RunnerID        uint64
	Price           float32
	Stake           float32
	Side            string
}

type Trade struct {
	Exchange  string
	Reference string
	Timestamp string
}
