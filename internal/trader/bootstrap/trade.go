package bootstrap

import (
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/postgres"
)

func (c Container) TradeWriter() trader.TradeWriter {
	return postgres.NewTradeWriter(c.Database)
}
