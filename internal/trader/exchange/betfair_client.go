package exchange

import (
	"context"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-betfair-go-client"
	"math"
)

type exchangeClient struct {
	client betfair.Client
	clock  clockwork.Clock
}

func (e *exchangeClient) Account(ctx context.Context) (*Account, error) {
	a, err := e.client.AccountFunds(ctx)

	if err != nil {
		return nil, err
	}

	return &Account{
		Balance:       a.Balance,
		Exposure:      a.Exposure,
		ExposureLimit: a.ExposureLimit,
	}, nil
}

func (e *exchangeClient) PlaceTrade(ctx context.Context, t *TradeTicket) (*Trade, error) {
	req := buildPlaceOrderRequest(t)

	res, err := e.client.PlaceOrder(ctx, req)

	if err != nil {
		return nil, &ClientError{action: "place orders", err: err}
	}

	if len(res.InstructionReports) != 1 {
		return nil, &InvalidResponseError{message: "response does not contain expected instruction report"}
	}

	report := res.InstructionReports[0]

	if res.Status != Success {
		return nil, &OrderFailureError{
			marketID:  t.MarketID,
			runnerID:  t.RunnerID,
			status:    report.Status,
			errorCode: report.ErrorCode,
			stake:     t.Stake,
			price:     t.Price,
		}
	}

	if report.OrderStatus != ExecutionComplete {
		return nil, &UnmatchedError{
			marketID: t.MarketID,
			runnerID: t.RunnerID,
			status:   report.OrderStatus,
		}
	}

	trade := Trade{
		Exchange:  Betfair,
		Reference: report.BetID,
		Timestamp: report.PlacedDate,
	}

	return &trade, nil
}

func buildPlaceOrderRequest(t *TradeTicket) betfair.PlaceOrderRequest {
	o := betfair.LimitOrder{
		Size:            float32(math.Round(float64(t.Stake)*100)/100),
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

func NewBetFairExchangeClient(c betfair.Client) Client {
	return &exchangeClient{client: c}
}
