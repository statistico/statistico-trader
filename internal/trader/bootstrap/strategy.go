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

func (c Container) StrategyFilterMatcher() strategy.FilterMatcher {
	return strategy.NewFilterMatcher(
		c.DataServiceFixtureClient(),
		c.StrategyResultClassifier(),
		c.StrategyStatClassifier(),
	)
}

func (c Container) StrategyResultClassifier() strategy.ResultFilterClassifier {
	return strategy.NewResultFilterClassifier(c.DataServiceResultClient())
}

func (c Container) StrategyStatClassifier() strategy.StatFilterClassifier {
	return strategy.NewStatFilterClassifier(c.DataServiceResultClient())
}

func (c Container) StrategyFinder() strategy.Finder {
	return strategy.NewFinder(c.StrategyReader(), c.StrategyFilterMatcher(), c.Logger)
}

func (c Container) StrategyResultParser() strategy.ResultParser {
	return strategy.NewResultParser(c.DataServiceResultClient())
}

func (c Container) StrategyBuilder() strategy.Builder {
	return strategy.NewBuilder(
		c.StrategyFilterMatcher(),
		c.StrategyResultParser(),
		c.OddsWarehouseMarketClient(),
		c.Logger,
	)
}
