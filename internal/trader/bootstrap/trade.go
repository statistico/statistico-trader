package bootstrap

import (
	"github.com/statistico/statistico-trader/internal/trader/trade"
)

func (c Container) TradeWriter() trade.Writer {
	return trade.NewPostgresWriter(c.Database)
}

func (c Container) TradeReader() trade.Reader {
	return trade.NewPostgresReader(c.Database)
}

func (c Container) TradePlacer() trade.Placer {
	return trade.NewPlacer(c.TradeReader(), c.TradeWriter(), c.Clock)
}

func (c Container) TradeManager() trade.Manager {
	return trade.NewManager(c.ExchangeClientFactory(), c.UserService(), c.TradePlacer())
}
