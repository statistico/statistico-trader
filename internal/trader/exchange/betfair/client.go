package betfair

import (
	"context"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-betfair-go-client"
	"github.com/statistico/statistico-strategy/internal/trader/exchange"
)

const (
	Betfair = "BETFAIR"
	ExecutionComplete = "EXECUTION_COMPLETE"
	FillOrKill = "FILL_OR_KILL"
	PersistenceTypeLapse = "LAPSE"
	OrderTypeLimit = "LIMIT"
	Success = "SUCCESS"
)

type exchangeClient struct {
	client betfair.Client
	clock  clockwork.Clock
}

func (e *exchangeClient) Account(ctx context.Context) (*exchange.Account, error) {
	a, err := e.client.AccountFunds(ctx)

	if err != nil {
		return nil, err
	}

	return &exchange.Account{
		Balance:       a.Balance,
		Exposure:      a.Exposure,
		ExposureLimit: a.ExposureLimit,
	}, nil
}

func (e *exchangeClient) PlaceTrade(ctx context.Context, t *exchange.TradeTicket) (*exchange.Trade, error) {
	req := buildPlaceOrderRequest(t)

	res, err := e.client.PlaceOrder(ctx, req)

	if err != nil {
		return nil, &exchange.ClientError{Action: "place orders", Err: err}
	}

	if len(res.InstructionReports) != 1 {
		return nil, &exchange.InvalidResponseError{Message: "response does not contain expected instruction report"}
	}

	report := res.InstructionReports[0]

	if res.Status != Success {
		return nil, &exchange.OrderFailureError{
			MarketID:  t.MarketID,
			RunnerID:  t.RunnerID,
			Status:    report.Status,
			ErrorCode: report.ErrorCode,
		}
	}

	if report.OrderStatus != ExecutionComplete {
		return nil, &exchange.UnmatchedError{
			MarketID: t.MarketID,
			RunnerID: t.RunnerID,
			Status:   report.OrderStatus,
		}
	}

	trade := exchange.Trade{
		Exchange:  Betfair,
		Reference: report.BetID,
		Timestamp: report.PlacedDate,
	}

	return &trade, nil
}

func buildPlaceOrderRequest(t *exchange.TradeTicket) betfair.PlaceOrderRequest {
	o := betfair.LimitOrder{
		Size:            t.Stake,
		Price:           t.Price,
		PersistenceType: PersistenceTypeLapse,
		TimeInForce:     FillOrKill,
	}

	i := betfair.PlaceInstruction{
		OrderType:   OrderTypeLimit,
		SelectionID: t.RunnerID,
		Side:        t.Side,
		LimitOrder:  o,
	}

	return betfair.PlaceOrderRequest{
		MarketID:            t.MarketID,
		Instructions:        []betfair.PlaceInstruction{i},
	}
}

func NewExchangeClient(c betfair.Client) exchange.Client {
	return &exchangeClient{client: c}
}
