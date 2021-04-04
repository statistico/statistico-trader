package grpc

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
)

type TradeQuery struct {
	Markets       <-chan *statistico.MarketRunner
	ResultFilters []*trader.ResultFilter
	StatFilters   []*trader.StatFilter
}
