package grpc

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader"
)

type TradeQuery struct {
	Markets       <-chan *statistico.MarketRunner
	ResultFilters []*trader.ResultFilter
	StatFilters   []*trader.StatFilter
}
