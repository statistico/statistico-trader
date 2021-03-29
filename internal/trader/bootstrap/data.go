package bootstrap

import statisticodata "github.com/statistico/statistico-data-go-grpc-client"

func (c Container) DataServiceResultClient() statisticodata.ResultClient {
	return statisticodata.NewResultClient(c.GrpcResultClient())
}

func (c Container) DataServiceFixtureClient() statisticodata.FixtureClient {
	return statisticodata.NewFixtureClient(c.GrpcFixtureClient())
}
