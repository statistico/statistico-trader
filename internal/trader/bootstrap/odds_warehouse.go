package bootstrap

import "github.com/statistico/statistico-odds-warehouse-go-grpc-client"

func (c Container) OddsWarehouseMarketClient() statisticooddswarehouse.MarketClient {
	return statisticooddswarehouse.NewMarketClient(c.GrpcMarketClient())
}
