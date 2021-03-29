package bootstrap

import "github.com/statistico/statistico-strategy/internal/trader/market"

func (c Container) TradeFactory() market.TradeFactory {
	return market.NewTradeFactory(
		c.DataServiceFixtureClient(),
		c.DataServiceResultClient(),
		c.FilterMatcher(),
		c.Clock,
	)
}
