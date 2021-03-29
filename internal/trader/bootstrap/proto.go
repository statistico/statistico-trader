package bootstrap

import (
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc"
)

func (c Container) GrpcFixtureClient() statistico.FixtureServiceClient {
	config := c.Config

	address := config.StatisticoDataService.Host + ":" + config.StatisticoDataService.Port

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		c.Logger.Warnf("Error initializing statistico data service grpc client %s", err.Error())
	}

	return statistico.NewFixtureServiceClient(conn)
}

func (c Container) GrpcResultClient() statistico.ResultServiceClient {
	config := c.Config

	address := config.StatisticoDataService.Host + ":" + config.StatisticoDataService.Port

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		c.Logger.Warnf("Error initializing statistico data service grpc client %s", err.Error())
	}

	return statistico.NewResultServiceClient(conn)
}

func (c Container) GrpcMarketClient() statistico.OddsWarehouseServiceClient {
	config := c.Config

	address := config.StatisticoOddsWarehouseService.Host + ":" + config.StatisticoOddsWarehouseService.Port

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		c.Logger.Warnf("Error initializing statistico data service grpc client %s", err.Error())
	}

	return statistico.NewOddsWarehouseServiceClient(conn)
}
