package exchange

import (
	bfc "github.com/statistico/statistico-betfair-go-client"
	"net/http"
)

type ClientFactory interface {
	Create(exchange, user, password, key string) (Client, error)
}

type clientFactory struct {
	httpClient *http.Client
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

	client := bfc.NewClient(c.httpClient, credentials)

	return NewBetFairExchangeClient(client), nil
}

func NewClientFactory(c *http.Client) ClientFactory {
	return &clientFactory{httpClient: c}
}
