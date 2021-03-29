package grpc_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/grpc"
	"github.com/statistico/statistico-strategy/internal/trader/market"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTradeFinder_Find(t *testing.T) {
	t.Run("returns a channel of statistico StrategyTrade struct", func(t *testing.T) {
		t.Helper()

		factory := new(MockTradeFactory)
		logger, _ := test.NewNullLogger()

		finder := grpc.NewTradeFinder(factory, logger)

		date, _ := ptypes.TimestampProto(time.Unix(1584014400, 0))

		markets := make([]*statistico.MarketRunner, 2)

		markets[0] = &statistico.MarketRunner{
			MarketName:    "OVER_UNDER_25",
			RunnerName:    "Over 2.5 Goals",
			EventId:       138171,
			CompetitionId: 8,
			SeasonId:      17420,
			EventDate:     date,
			Exchange:      "betfair",
			Price: &statistico.Price{
				Value: 1.95,
				Side:  statistico.SideEnum_BACK,
			},
		}

		markets[1] = &statistico.MarketRunner{
			MarketName:    "OVER_UNDER_25",
			RunnerName:    "Over 2.5 Goals",
			EventId:       138172,
			CompetitionId: 8,
			SeasonId:      17420,
			EventDate:     date,
			Exchange:      "betfair",
			Price: &statistico.Price{
				Value: 2.15,
				Side:  statistico.SideEnum_BACK,
			},
		}

		query := &grpc.TradeQuery{
			Markets: markerRunnerChannel(markets),
			RunnerFilters: []*trader.ResultFilter{
				{
					Team:   "HOME_TEAM",
					Result: "WIN",
					Games:  3,
					Venue:  "HOME",
				},
			},
		}

		marketQuery := mock.MatchedBy(func(q *market.Query) bool {
			a := assert.New(t)

			if q.EventId == 138171 {
				a.Equal("OVER_UNDER_25", q.MarketName)
				a.Equal("Over 2.5 Goals", q.RunnerName)
				a.Equal(uint64(138171), q.EventId)
				a.Equal(uint64(8), q.CompetitionId)
				a.Equal(uint64(17420), q.SeasonId)
				a.Equal("2020-03-12T12:00:00Z", q.EventDate.Format(time.RFC3339))
				a.Equal("BACK", q.Side)
				a.Equal("betfair", q.Exchange)
				a.Equal(float32(1.95), q.RunnerPrice)
				a.Equal(query.RunnerFilters, q.ResultFilters)
				a.Equal(query.StatFilters, q.StatFilters)
			}

			if q.EventId == 138172 {
				a.Equal("OVER_UNDER_25", q.MarketName)
				a.Equal("Over 2.5 Goals", q.RunnerName)
				a.Equal(uint64(138172), q.EventId)
				a.Equal(uint64(8), q.CompetitionId)
				a.Equal(uint64(17420), q.SeasonId)
				a.Equal("2020-03-12T12:00:00Z", q.EventDate.Format(time.RFC3339))
				a.Equal("BACK", q.Side)
				a.Equal("betfair", q.Exchange)
				a.Equal(float32(2.15), q.RunnerPrice)
				a.Equal(query.RunnerFilters, q.ResultFilters)
				a.Equal(query.StatFilters, q.StatFilters)
			}

			return true
		})

		success := "SUCCESS"

		tradeOne := &market.Trade{
			MarketName:    "OVER_UNDER_25",
			RunnerName:    "Over 2.5 Goals",
			RunnerPrice:   1.95,
			EventId:       138171,
			CompetitionId: 8,
			SeasonId:      17420,
			EventDate:     time.Unix(1584014400, 0),
			Side:          "BACK",
			Result:        &success,
		}

		tradeTwo := &market.Trade{
			MarketName:    "OVER_UNDER_25",
			RunnerName:    "Over 2.5 Goals",
			RunnerPrice:   2.15,
			EventId:       138172,
			CompetitionId: 8,
			SeasonId:      17420,
			EventDate:     time.Unix(1584014400, 0),
			Side:          "BACK",
			Result:        &success,
		}

		ctx := context.Background()

		factory.On("CreateTrade", ctx, marketQuery).Return(tradeOne, nil).Once()
		factory.On("CreateTrade", ctx, marketQuery).Return(tradeTwo, nil).Once()

		tradeCh := finder.Find(ctx, query)

		one := <-tradeCh
		two := <-tradeCh

		a := assert.New(t)
		a.Equal(uint64(138171), one.EventId)
		a.Equal(uint64(138172), two.EventId)
		factory.AssertExpectations(t)
	})

	t.Run("logs an error if returned by TradeFactory", func(t *testing.T) {
		t.Helper()

		factory := new(MockTradeFactory)
		logger, hook := test.NewNullLogger()

		finder := grpc.NewTradeFinder(factory, logger)

		date, _ := ptypes.TimestampProto(time.Unix(1584014400, 0))

		markets := make([]*statistico.MarketRunner, 1)

		markets[0] = &statistico.MarketRunner{
			MarketName:    "OVER_UNDER_25",
			RunnerName:    "Over 2.5 Goals",
			EventId:       138171,
			CompetitionId: 8,
			SeasonId:      17420,
			EventDate:     date,
			Exchange:      "betfair",
			Price: &statistico.Price{
				Value: 1.95,
				Side:  statistico.SideEnum_BACK,
			},
		}

		query := &grpc.TradeQuery{
			Markets: markerRunnerChannel(markets),
			RunnerFilters: []*trader.ResultFilter{
				{
					Team:   "HOME_TEAM",
					Result: "WIN",
					Games:  3,
					Venue:  "HOME",
				},
			},
		}

		marketQuery := mock.MatchedBy(func(q *market.Query) bool {
			a := assert.New(t)

			a.Equal("OVER_UNDER_25", q.MarketName)
			a.Equal("Over 2.5 Goals", q.RunnerName)
			a.Equal(uint64(138171), q.EventId)
			a.Equal(uint64(8), q.CompetitionId)
			a.Equal(uint64(17420), q.SeasonId)
			a.Equal("2020-03-12T12:00:00Z", q.EventDate.Format(time.RFC3339))
			a.Equal("BACK", q.Side)
			a.Equal("betfair", q.Exchange)
			a.Equal(float32(1.95), q.RunnerPrice)
			a.Equal(query.RunnerFilters, q.ResultFilters)
			a.Equal(query.StatFilters, q.StatFilters)

			return true
		})

		success := "SUCCESS"

		tradeOne := &market.Trade{
			MarketName:    "OVER_UNDER_25",
			RunnerName:    "Over 2.5 Goals",
			RunnerPrice:   1.95,
			EventId:       138171,
			CompetitionId: 8,
			SeasonId:      17420,
			EventDate:     time.Unix(1584014400, 0),
			Side:          "BACK",
			Result:        &success,
		}

		ctx := context.Background()

		factory.On("CreateTrade", ctx, marketQuery).Return(tradeOne, errors.New("filter error")).Once()

		tradeCh := finder.Find(ctx, query)

		trade := <-tradeCh

		a := assert.New(t)
		a.Nil(trade)
		a.Equal(1, len(hook.AllEntries()))
		a.Equal("error handling trade for market OVER_UNDER_25, runner Over 2.5 Goals and event 138171: filter error", hook.LastEntry().Message)
		factory.AssertExpectations(t)
	})
}

func markerRunnerChannel(mr []*statistico.MarketRunner) chan *statistico.MarketRunner {
	ch := make(chan *statistico.MarketRunner, len(mr))

	for _, m := range mr {
		ch <- m
	}

	close(ch)

	return ch
}

type MockTradeFactory struct {
	mock.Mock
}

func (m *MockTradeFactory) CreateTrade(ctx context.Context, q *market.Query) (*market.Trade, error) {
	args := m.Called(ctx, q)
	return args.Get(0).(*market.Trade), args.Error(1)
}
