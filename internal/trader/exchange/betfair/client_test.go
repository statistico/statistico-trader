package betfair_test

import (
	"context"
	"errors"
	betfair2 "github.com/statistico/statistico-betfair-go-client"
	"github.com/statistico/statistico-strategy/internal/trader/exchange"
	"github.com/statistico/statistico-strategy/internal/trader/exchange/betfair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestExchangeClient_Account(t *testing.T) {
	t.Run("returns account information for a user", func(t *testing.T) {
		t.Helper()

		mc := new(betfair2.MockClient)

		client := betfair.NewExchangeClient(mc)

		af := betfair2.AccountFunds{
			Balance:            10.00,
			DiscountRate:       0,
			Exposure:           -59.82,
			ExposureLimit:      -5000.00,
			PointBalance:       0,
			RetainedCommission: 0,
			Wallet:             "WALLET",
		}

		ctx := context.Background()

		mc.On("AccountFunds", ctx).Return(&af, nil)

		account, err := client.Account(ctx)

		if err != nil {
			t.Fatalf("Expected nil, got %+v", err)
		}

		a := assert.New(t)

		a.Equal(float32(10.0), account.Balance)
		a.Equal(float32(-59.82), account.Exposure)
		a.Equal(float32(-5000.00), account.ExposureLimit)
	})

	t.Run("returns an error if error returned by betfair client", func(t *testing.T) {
		t.Helper()

		mc := new(betfair2.MockClient)

		client := betfair.NewExchangeClient(mc)

		ctx := context.Background()

		mc.On("AccountFunds", ctx).Return(&betfair2.AccountFunds{}, errors.New("client error"))

		_, err := client.Account(ctx)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "client error", err.Error())
	})
}

func TestExchangeClient_PlaceTrade(t *testing.T) {
	t.Run("places trade via betfair client and returns exchange.Trade struct", func(t *testing.T) {
		t.Helper()

		mc := new(betfair2.MockClient)
		client := betfair.NewExchangeClient(mc)

		ticket := exchange.TradeTicket{
			MarketID:        "1.181098580",
			RunnerID:        16082847,
			Price:           19.0,
			Stake:           2.0,
			Side:            "BACK",
		}

		req := mock.MatchedBy(func(r betfair2.PlaceOrderRequest) bool {
			o := betfair2.LimitOrder{
				Size:            2.0,
				Price:           19.0,
				PersistenceType: "LAPSE",
				TimeInForce:     "FILL_OR_KILL",
			}

			i := betfair2.PlaceInstruction{
				OrderType:   "LIMIT",
				SelectionID: 16082847,
				Side:        "BACK",
				LimitOrder:  o,
			}

			por := betfair2.PlaceOrderRequest{
				MarketID:            "1.181098580",
				Instructions:        []betfair2.PlaceInstruction{i},
				CustomerRef:         "",
				MarketVersion:       0,
			}
			
			assert.Equal(t, por, r)

			return true
		})
		
		res := betfair2.PlaceExecutionReport{
			MarketID:           "1.181098580",
			Status:             "SUCCESS",
			InstructionReports: []betfair2.PlaceInstructionReport{
				{
					Status:              "SUCCESS",
					OrderStatus:         "EXECUTION_COMPLETE",
					Instruction:         betfair2.PlaceInstruction{},
					BetID:               "BET-ID-123",
					PlacedDate:          "2020-04-07T12:00:00+00:00",
					SizeMatched:         19.0,
					AveragePriceMatched: 0,
				},
			},
		}

		ctx := context.Background()

		mc.On("PlaceOrder", ctx, req).Return(&res, nil)

		trade, err := client.PlaceTrade(ctx, &ticket)

		if err != nil {
			t.Fatalf("Expected nil, got %+v", err)
		}

		a := assert.New(t)

		a.Equal("BETFAIR", trade.Exchange)
		a.Equal("BET-ID-123", trade.Reference)
		a.Equal("2020-04-07T12:00:00+00:00", trade.Timestamp)
	})

	t.Run("returns exchange.ClientError if error returned by betfair client", func(t *testing.T) {
		t.Helper()

		mc := new(betfair2.MockClient)
		client := betfair.NewExchangeClient(mc)

		ticket := exchange.TradeTicket{
			MarketID:        "1.181098580",
			RunnerID:        16082847,
			Price:           19.0,
			Stake:           2.0,
			Side:            "BACK",
		}

		req := mock.MatchedBy(func(r betfair2.PlaceOrderRequest) bool {
			o := betfair2.LimitOrder{
				Size:            2.0,
				Price:           19.0,
				PersistenceType: "LAPSE",
				TimeInForce:     "FILL_OR_KILL",
			}

			i := betfair2.PlaceInstruction{
				OrderType:   "LIMIT",
				SelectionID: 16082847,
				Side:        "BACK",
				LimitOrder:  o,
			}

			por := betfair2.PlaceOrderRequest{
				MarketID:            "1.181098580",
				Instructions:        []betfair2.PlaceInstruction{i},
				CustomerRef:         "",
				MarketVersion:       0,
			}

			assert.Equal(t, por, r)

			return true
		})

		e := errors.New("client error")

		ctx := context.Background()

		mc.On("PlaceOrder", ctx, req).Return(&betfair2.PlaceExecutionReport{}, e)

		_, err := client.PlaceTrade(ctx, &ticket)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "error making 'place orders' request: client error", err.Error())
	})

	t.Run("returns exchange.InvalidResponseError if response does not contain InstructionReport", func(t *testing.T) {
		t.Helper()

		mc := new(betfair2.MockClient)
		client := betfair.NewExchangeClient(mc)

		ticket := exchange.TradeTicket{
			MarketID:        "1.181098580",
			RunnerID:        16082847,
			Price:           19.0,
			Stake:           2.0,
			Side:            "BACK",
		}

		req := mock.MatchedBy(func(r betfair2.PlaceOrderRequest) bool {
			o := betfair2.LimitOrder{
				Size:            2.0,
				Price:           19.0,
				PersistenceType: "LAPSE",
				TimeInForce:     "FILL_OR_KILL",
			}

			i := betfair2.PlaceInstruction{
				OrderType:   "LIMIT",
				SelectionID: 16082847,
				Side:        "BACK",
				LimitOrder:  o,
			}

			por := betfair2.PlaceOrderRequest{
				MarketID:            "1.181098580",
				Instructions:        []betfair2.PlaceInstruction{i},
				CustomerRef:         "",
				MarketVersion:       0,
			}

			assert.Equal(t, por, r)

			return true
		})

		res := betfair2.PlaceExecutionReport{
			MarketID:           "1.181098580",
			Status:             "SUCCESS",
		}

		ctx := context.Background()

		mc.On("PlaceOrder", ctx, req).Return(&res, nil)

		_, err := client.PlaceTrade(ctx, &ticket)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "invalid response in exchange client: response does not contain expected instruction report", err.Error())
	})

	t.Run("returns exchange.OrderFailureError if order status is not SUCCESS", func(t *testing.T) {
		t.Helper()

		mc := new(betfair2.MockClient)
		client := betfair.NewExchangeClient(mc)

		ticket := exchange.TradeTicket{
			MarketID:        "1.181098580",
			RunnerID:        16082847,
			Price:           19.0,
			Stake:           2.0,
			Side:            "BACK",
		}

		req := mock.MatchedBy(func(r betfair2.PlaceOrderRequest) bool {
			o := betfair2.LimitOrder{
				Size:            2.0,
				Price:           19.0,
				PersistenceType: "LAPSE",
				TimeInForce:     "FILL_OR_KILL",
			}

			i := betfair2.PlaceInstruction{
				OrderType:   "LIMIT",
				SelectionID: 16082847,
				Side:        "BACK",
				LimitOrder:  o,
			}

			por := betfair2.PlaceOrderRequest{
				MarketID:            "1.181098580",
				Instructions:        []betfair2.PlaceInstruction{i},
				CustomerRef:         "",
				MarketVersion:       0,
			}

			assert.Equal(t, por, r)

			return true
		})

		res := betfair2.PlaceExecutionReport{
			MarketID:           "1.181098580",
			Status:             "FAILURE",
			InstructionReports: []betfair2.PlaceInstructionReport{
				{
					Status:              "FAILURE",
					OrderStatus:         "EXPIRED",
					ErrorCode:           "INVALID_MARKET_ID",
					Instruction:         betfair2.PlaceInstruction{},
					BetID:               "BET-ID-123",
					PlacedDate:          "2020-04-07T12:00:00+00:00",
					SizeMatched:         19.0,
					AveragePriceMatched: 0,
				},
			},
		}

		ctx := context.Background()

		mc.On("PlaceOrder", ctx, req).Return(&res, nil)

		_, err := client.PlaceTrade(ctx, &ticket)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(
			t,
			"error placing order for market '1.181098580' and runner '16082847'. Code: 'INVALID_MARKET_ID' and Status: 'FAILURE'",
			err.Error(),
		)
	})

	t.Run("returns exchange.UnmatchedError if report order status is not EXECUTION_COMPLETE", func(t *testing.T) {
		t.Helper()

		mc := new(betfair2.MockClient)
		client := betfair.NewExchangeClient(mc)

		ticket := exchange.TradeTicket{
			MarketID:        "1.181098580",
			RunnerID:        16082847,
			Price:           19.0,
			Stake:           2.0,
			Side:            "BACK",
		}

		req := mock.MatchedBy(func(r betfair2.PlaceOrderRequest) bool {
			o := betfair2.LimitOrder{
				Size:            2.0,
				Price:           19.0,
				PersistenceType: "LAPSE",
				TimeInForce:     "FILL_OR_KILL",
			}

			i := betfair2.PlaceInstruction{
				OrderType:   "LIMIT",
				SelectionID: 16082847,
				Side:        "BACK",
				LimitOrder:  o,
			}

			por := betfair2.PlaceOrderRequest{
				MarketID:            "1.181098580",
				Instructions:        []betfair2.PlaceInstruction{i},
				CustomerRef:         "",
				MarketVersion:       0,
			}

			assert.Equal(t, por, r)

			return true
		})

		res := betfair2.PlaceExecutionReport{
			MarketID:           "1.181098580",
			Status:             "SUCCESS",
			InstructionReports: []betfair2.PlaceInstructionReport{
				{
					Status:              "SUCCESS",
					OrderStatus:         "EXPIRED",
					Instruction:         betfair2.PlaceInstruction{},
					BetID:               "BET-ID-123",
					PlacedDate:          "2020-04-07T12:00:00+00:00",
					SizeMatched:         19.0,
					AveragePriceMatched: 0,
				},
			},
		}

		ctx := context.Background()

		mc.On("PlaceOrder", ctx, req).Return(&res, nil)

		_, err := client.PlaceTrade(ctx, &ticket)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "trade unmatched for market '1.181098580' and runner '16082847' with status 'EXPIRED'", err.Error())
	})
}
