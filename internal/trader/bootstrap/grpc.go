package bootstrap

import "github.com/statistico/statistico-trader/internal/trader/grpc"

func (c Container) GrpcStrategyService() *grpc.StrategyService {
	return grpc.NewStrategyService(
		c.StrategyWriter(),
		c.StrategyReader(),
		c.OddsWarehouseMarketClient(),
		c.Logger,
		c.Clock,
	)
}
