package grpc

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
)

type TradeQuery struct {
	Markets       <-chan *statistico.MarketRunner
	ResultFilters []*strategy.ResultFilter
	StatFilters   []*strategy.StatFilter
}
