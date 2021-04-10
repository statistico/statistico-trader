package bootstrap

import (
	"github.com/statistico/statistico-trader/internal/trader/strategy"
)

func (c Container) StrategyWriter() strategy.Writer {
	return strategy.NewPostgresWriter(c.Database)
}

func (c Container) StrategyReader() strategy.Reader {
	return strategy.NewPostgresReader(c.Database)
}
