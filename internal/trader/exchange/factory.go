package exchange

import (
	bfc "github.com/statistico/statistico-betfair-go-client"
	"github.com/statistico/statistico-trader/internal/trader/bootstrap"
	"github.com/statistico/statistico-trader/internal/trader/exchange/betfair"
)

const (
	Betfair = "betfair"
)

type ClientFactory interface {
	Create(exchange, user, password, key string) (Client, error)
}

type clientFactory struct {
	config *bootstrap.Config
}

func (c *clientFactory) Create(exchange, user, password, key string) (Client, error) {
	if exchange != Betfair {
		return nil, &InvalidExchangeError{exchange: exchange}
	}

	credentials := bfc.InteractiveCredentials{
		Username: user,
		Password: password,
		Key:      key,
	}

	client := bfc.NewClient(c.config.HTTPClient, credentials)

	return betfair.NewExchangeClient(client), nil
}

func NewClientFactory(c *bootstrap.Config) ClientFactory {
	return &clientFactory{config: c}
}
