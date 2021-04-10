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
