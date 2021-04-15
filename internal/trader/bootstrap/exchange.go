package bootstrap

import "github.com/statistico/statistico-trader/internal/trader/exchange"

func (c Container) ExchangeClientFactory() exchange.ClientFactory {
	return exchange.NewClientFactory(c.Config.HTTPClient)
}
