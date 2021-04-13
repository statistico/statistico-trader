package trade_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/statistico/statistico-trader/internal/trader/trade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestPlacer_PlaceTrade(t *testing.T) {
	ticket := trade.Ticket{
		MarketID:      "1.234567",
		MarketName:    "MATCH_ODDS",
		RunnerID:      9876,
		RunnerName:    "Home",
		EventID:       345192,
		CompetitionID: 8,
		SeasonID:      17420,
		EventDate:     time.Unix(1615723200, 0),
		Exchange:      "betfair",
		Price:         trade.TicketPrice{
			Value: 1.94,
			Size:  500.00,
			Side:  "BACK",
		},
	}

	st := strategy.Strategy{
		ID:             uuid.MustParse("9dbc01ae-bea0-45b7-a3b1-92ae095dfad0"),
		Name:           "Joe's Super Strategy",
		Description:    "The strategy to make me a millionaire",
		UserID:         uuid.New(),
		MarketName:     "MATCH_ODDS",
		RunnerName:     "Home",
		StakingPlan:    strategy.StakingPlan{
			Name:   "PERCENTAGE",
			Number: 10,
		},
	}

	t.Run("uses exchange.Client to place trade and inserts via trade.Writer", func(t *testing.T) {
		t.Helper()

		reader := new(MockTradeReader)
		writer := new(MockTradeWriter)
		clock := clockwork.NewFakeClockAt(time.Unix(1615550400, 0))
		placer := trade.NewPlacer(reader, writer, clock)

		ctx := context.Background()
		client := new(MockExchangeClient)

		account := exchange.Account{
			Balance:       500,
			Exposure:      -100,
			ExposureLimit: -5000,
		}

		reader.On("Exists", ticket.MarketName, ticket.RunnerName, ticket.EventID, st.ID).Return(false, nil)

		client.On("Account", ctx).Return(&account, nil)

		mockTicket := mock.MatchedBy(func(e *exchange.TradeTicket) bool {
			a := assert.New(t)

			a.Equal("1.234567", e.MarketID)
			a.Equal(uint64(9876), e.RunnerID)
			a.Equal(float32(1.94), e.Price)
			a.Equal(float32(60.00), e.Stake)
			a.Equal("BACK", e.Side)
			return true
		})

		response := exchange.Trade{
			Exchange:  "betfair",
			Reference: "1234567890",
		}

		client.On("PlaceTrade", ctx, mockTicket).Return(&response, nil)

		mockTrade := mock.MatchedBy(func(tr *trade.Trade) bool {
			a := assert.New(t)

			a.Equal(st.ID, tr.StrategyID)
			a.Equal("betfair", tr.Exchange)
			a.Equal("1234567890", tr.ExchangeRef)
			a.Equal(ticket.MarketName, tr.Market)
			a.Equal(ticket.RunnerName, tr.Runner)
			a.Equal(float32(1.94), tr.Price)
			a.Equal(float32(60.00), tr.Stake)
			a.Equal(ticket.EventID, tr.EventID)
			a.Equal(ticket.EventDate, tr.EventDate)
			a.Equal("BACK", tr.Side)
			a.Equal("IN_PLAY", tr.Result)
			a.Equal(clock.Now(), tr.Timestamp)
			return true
		})

		writer.On("Insert", mockTrade).Return(nil)

		tr, err := placer.PlaceTrade(ctx, client, &ticket, &st)

		if err != nil {
			t.Fatalf("Expected nil, got %+v", err)
		}

		a := assert.New(t)

		a.Equal(st.ID, tr.StrategyID)
		a.Equal("betfair", tr.Exchange)
		a.Equal("1234567890", tr.ExchangeRef)
		a.Equal(ticket.MarketName, tr.Market)
		a.Equal(ticket.RunnerName, tr.Runner)
		a.Equal(float32(1.94), tr.Price)
		a.Equal(float32(60.00), tr.Stake)
		a.Equal(ticket.EventID, tr.EventID)
		a.Equal(ticket.EventDate, tr.EventDate)
		a.Equal("BACK", tr.Side)
		a.Equal("IN_PLAY", tr.Result)
		a.Equal(clock.Now(), tr.Timestamp)

		reader.AssertExpectations(t)
		writer.AssertExpectations(t)
		client.AssertExpectations(t)
	})

	t.Run("returns a DuplicationError is trade already exists", func(t *testing.T) {
		t.Helper()

		reader := new(MockTradeReader)
		writer := new(MockTradeWriter)
		clock := clockwork.NewFakeClockAt(time.Unix(1615550400, 0))
		placer := trade.NewPlacer(reader, writer, clock)

		ctx := context.Background()
		client := new(MockExchangeClient)

		reader.On("Exists", ticket.MarketName, ticket.RunnerName, ticket.EventID, st.ID).Return(true, nil)

		client.AssertNotCalled(t, "Account")
		client.AssertNotCalled(t, "PlaceTrade")
		writer.AssertNotCalled(t, "Insert")

		_, err := placer.PlaceTrade(ctx, client, &ticket, &st)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "trade exists for market MATCH_ODDS, runner Home, event 345192 and strategy 9dbc01ae-bea0-45b7-a3b1-92ae095dfad0", err.Error())
	})

	t.Run("returns an ExchangeError if error is returned by exchange.Client", func(t *testing.T) {
		t.Helper()

		reader := new(MockTradeReader)
		writer := new(MockTradeWriter)
		clock := clockwork.NewFakeClockAt(time.Unix(1615550400, 0))
		placer := trade.NewPlacer(reader, writer, clock)

		ctx := context.Background()
		client := new(MockExchangeClient)

		reader.On("Exists", ticket.MarketName, ticket.RunnerName, ticket.EventID, st.ID).Return(false, nil)

		client.On("Account", ctx).Return(&exchange.Account{}, errors.New("client error"))

		client.AssertNotCalled(t, "PlaceTrade")
		writer.AssertNotCalled(t, "Insert")

		_, err := placer.PlaceTrade(ctx, client, &ticket, &st)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "error returned by exchange client: client error", err.Error())
	})

	t.Run("returns InvalidBalanceError if exchange balance is zero or less", func(t *testing.T) {
		t.Helper()

		reader := new(MockTradeReader)
		writer := new(MockTradeWriter)
		clock := clockwork.NewFakeClockAt(time.Unix(1615550400, 0))
		placer := trade.NewPlacer(reader, writer, clock)

		ctx := context.Background()
		client := new(MockExchangeClient)

		account := exchange.Account{
			Balance:       0,
			Exposure:      -10,
			ExposureLimit: -5000,
		}

		reader.On("Exists", ticket.MarketName, ticket.RunnerName, ticket.EventID, st.ID).Return(false, nil)

		client.On("Account", ctx).Return(&account, nil)

		client.AssertNotCalled(t, "PlaceTrade")
		writer.AssertNotCalled(t, "Insert")

		_, err := placer.PlaceTrade(ctx, client, &ticket, &st)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "invalid balance of 0.00 when placing trade for market MATCH_ODDS, runner Home, event 345192 and strategy 9dbc01ae-bea0-45b7-a3b1-92ae095dfad0", err.Error())
	})

	t.Run("returns an ExchangeError if error returned placing a trade", func(t *testing.T) {
		t.Helper()

		reader := new(MockTradeReader)
		writer := new(MockTradeWriter)
		clock := clockwork.NewFakeClockAt(time.Unix(1615550400, 0))
		placer := trade.NewPlacer(reader, writer, clock)

		ctx := context.Background()
		client := new(MockExchangeClient)

		account := exchange.Account{
			Balance:       500,
			Exposure:      -100,
			ExposureLimit: -5000,
		}

		reader.On("Exists", ticket.MarketName, ticket.RunnerName, ticket.EventID, st.ID).Return(false, nil)

		client.On("Account", ctx).Return(&account, nil)

		mockTicket := mock.MatchedBy(func(e *exchange.TradeTicket) bool {
			a := assert.New(t)

			a.Equal("1.234567", e.MarketID)
			a.Equal(uint64(9876), e.RunnerID)
			a.Equal(float32(1.94), e.Price)
			a.Equal(float32(60.00), e.Stake)
			a.Equal("BACK", e.Side)
			return true
		})

		client.On("PlaceTrade", ctx, mockTicket).Return(&exchange.Trade{}, errors.New("client error"))

		writer.AssertNotCalled(t, "Insert")

		_, err := placer.PlaceTrade(ctx, client, &ticket, &st)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "error returned by exchange client: client error", err.Error())
	})

	t.Run("returns trade and error if error inserting trade via trade.Writer", func(t *testing.T) {
		t.Helper()

		reader := new(MockTradeReader)
		writer := new(MockTradeWriter)
		clock := clockwork.NewFakeClockAt(time.Unix(1615550400, 0))
		placer := trade.NewPlacer(reader, writer, clock)

		ctx := context.Background()
		client := new(MockExchangeClient)

		account := exchange.Account{
			Balance:       500,
			Exposure:      -100,
			ExposureLimit: -5000,
		}

		reader.On("Exists", ticket.MarketName, ticket.RunnerName, ticket.EventID, st.ID).Return(false, nil)

		client.On("Account", ctx).Return(&account, nil)

		mockTicket := mock.MatchedBy(func(e *exchange.TradeTicket) bool {
			a := assert.New(t)

			a.Equal("1.234567", e.MarketID)
			a.Equal(uint64(9876), e.RunnerID)
			a.Equal(float32(1.94), e.Price)
			a.Equal(float32(60.00), e.Stake)
			a.Equal("BACK", e.Side)
			return true
		})

		response := exchange.Trade{
			Exchange:  "betfair",
			Reference: "1234567890",
		}

		client.On("PlaceTrade", ctx, mockTicket).Return(&response, nil)

		mockTrade := mock.MatchedBy(func(tr *trade.Trade) bool {
			a := assert.New(t)

			a.Equal(st.ID, tr.StrategyID)
			a.Equal("betfair", tr.Exchange)
			a.Equal("1234567890", tr.ExchangeRef)
			a.Equal(ticket.MarketName, tr.Market)
			a.Equal(ticket.RunnerName, tr.Runner)
			a.Equal(float32(1.94), tr.Price)
			a.Equal(float32(60.00), tr.Stake)
			a.Equal(ticket.EventID, tr.EventID)
			a.Equal(ticket.EventDate, tr.EventDate)
			a.Equal("BACK", tr.Side)
			a.Equal("IN_PLAY", tr.Result)
			a.Equal(clock.Now(), tr.Timestamp)
			return true
		})

		writer.On("Insert", mockTrade).Return(errors.New("error inserting trade"))

		tr, err := placer.PlaceTrade(ctx, client, &ticket, &st)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if tr == nil {
			t.Fatal("Expected trade, got nil")
		}

		a := assert.New(t)

		a.Equal(st.ID, tr.StrategyID)
		a.Equal("betfair", tr.Exchange)
		a.Equal("1234567890", tr.ExchangeRef)
		a.Equal(ticket.MarketName, tr.Market)
		a.Equal(ticket.RunnerName, tr.Runner)
		a.Equal(float32(1.94), tr.Price)
		a.Equal(float32(60.00), tr.Stake)
		a.Equal(ticket.EventID, tr.EventID)
		a.Equal(ticket.EventDate, tr.EventDate)
		a.Equal("BACK", tr.Side)
		a.Equal("IN_PLAY", tr.Result)
		a.Equal(clock.Now(), tr.Timestamp)
		a.Equal("error inserting trade", err.Error())

		reader.AssertExpectations(t)
		writer.AssertExpectations(t)
		client.AssertExpectations(t)
	})
}

type MockTradeReader struct {
	mock.Mock
}

func (m *MockTradeReader) Get(q *trade.ReaderQuery) ([]*trade.Trade, error) {
	args := m.Called(q)
	return args.Get(0).([]*trade.Trade), args.Error(1)
}

func (m *MockTradeReader) Exists(market, runner string, eventID uint64, strategyID uuid.UUID) (bool, error) {
	args := m.Called(market, runner, eventID, strategyID)
	return args.Get(0).(bool), args.Error(1)
}

type MockTradeWriter struct {
	mock.Mock
}

func (m *MockTradeWriter) Insert(t *trade.Trade) error {
	args := m.Called(t)
	return args.Error(0)
}

type MockExchangeClient struct {
	mock.Mock
}

func (m *MockExchangeClient) Account(ctx context.Context) (*exchange.Account, error) {
	args := m.Called(ctx)
	return args.Get(0).(*exchange.Account), args.Error(1)
}

func (m *MockExchangeClient) PlaceTrade(ctx context.Context, t *exchange.TradeTicket) (*exchange.Trade, error) {
	args := m.Called(ctx, t)
	return args.Get(0).(*exchange.Trade), args.Error(1)
}
