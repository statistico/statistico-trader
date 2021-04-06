package betfair

import (
	"github.com/statistico/statistico-betfair-go-client"
	"github.com/statistico/statistico-strategy/internal/trader/exchange"
)

type exchangeClient struct {
	client  betfair.Client
}

func (e *exchangeClient) Balance() (float32, error) {
	return 100000.54, nil
}

func (e *exchangeClient) PlaceTrade(t *exchange.TradeTicket) (*exchange.Trade, error) {
	return &exchange.Trade{}, nil
}

func NewExchangeClient(c betfair.Client) exchange.Client {
	return &exchangeClient{client: c}
}
