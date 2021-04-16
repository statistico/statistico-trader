package bootstrap

import "github.com/statistico/statistico-trader/internal/trader/grpc"

func (c Container) GrpcStrategyService() *grpc.StrategyService {
	return grpc.NewStrategyService(
		c.StrategyBuilder(),
		c.StrategyWriter(),
		c.StrategyReader(),
		c.Logger,
		c.Clock,
	)
}
