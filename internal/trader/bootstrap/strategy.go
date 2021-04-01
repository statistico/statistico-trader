package bootstrap

import (
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/postgres"
)

func (c Container) StrategyWriter() trader.StrategyWriter {
	return postgres.NewStrategyWriter(c.Database)
}

func (c Container) StrategyReader() trader.StrategyReader {
	return postgres.NewStrategyReader(c.Database)
}
