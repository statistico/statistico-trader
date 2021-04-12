package bootstrap

import "github.com/statistico/statistico-trader/internal/trader/market"

func (c Container) MarketHandler() market.Handler {
	return market.NewHandler(c.StrategyFinder(), c.TradeManager(), c.Logger)
}
