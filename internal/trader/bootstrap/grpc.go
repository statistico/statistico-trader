package bootstrap

import "github.com/statistico/statistico-strategy/internal/trader/grpc"

func (c Container) GrpcStrategyService() *grpc.StrategyService {
	return grpc.NewStrategyService(
		c.StrategyWriter(),
		c.OddsWarehouseMarketClient(),
		c.TradeFinder(),
		c.Logger,
		c.Clock,
	)
}

func (c Container) TradeFinder() grpc.TradeFinder {
	return grpc.NewTradeFinder(c.TradeFactory(), c.Logger)
}
