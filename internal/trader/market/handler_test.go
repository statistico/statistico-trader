package market_test

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-trader/internal/trader/market"
	"github.com/statistico/statistico-trader/internal/trader/queue"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/statistico/statistico-trader/internal/trader/trade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestHandler_HandleEventMarket(t *testing.T) {
	t.Run("finds matching strategies and manages trades for queue.EventMarket", func(t *testing.T) {
		t.Helper()

		finder := new(MockStrategyFinder)
		manager := new(MockTradeManager)
		logger, hook := test.NewNullLogger()
		handler := market.NewHandler(finder, manager, logger)

		event := queue.EventMarket{
			ID:            "1.234501",
			EventID:       45617,
			Name:          "OVER_UNDER_25",
			CompetitionID: 8,
			SeasonID:      17420,
			EventDate:     time.Unix(1616936636, 0),
			Exchange:      "betfair",
			Runners:       []*queue.Runner{
				{
					ID:         1,
					Name:       "Yes",
					Sort:       1,
					BackPrices: []queue.PriceSize{
						{
							Price: 1.95,
							Size:  500,
						},
					},
					LayPrices:  []queue.PriceSize{
						{
							Price: 1.99,
							Size:  100,
						},
					},
				},
				{
					ID:         2,
					Name:       "No",
					Sort:       1,
					BackPrices: []queue.PriceSize{
						{
							Price: 3.90,
							Size:  500,
						},
					},
					LayPrices:  []queue.PriceSize{
						{
							Price: 3.99,
							Size:  100,
						},
					},
				},
			},
			Timestamp:     1616936636,
		}

		ctx := context.Background()

		stOne := &strategy.Strategy{}

		finder.On("FindMatchingStrategies", ctx, mock.AnythingOfType("*strategy.FinderQuery")).
			Times(4).
			Return(strategyChannel(stOne))

		manager.On("Manage", ctx, mock.AnythingOfType("*trade.Ticket"), mock.AnythingOfType("*strategy.Strategy")).
			Times(4).
			Return(nil)

		handler.HandleEventMarket(ctx, &event)

		assert.Equal(t, 0, len(hook.AllEntries()))
	})

	t.Run("logs error if error returned by trade.Manager", func(t *testing.T) {
		t.Helper()

		finder := new(MockStrategyFinder)
		manager := new(MockTradeManager)
		logger, hook := test.NewNullLogger()
		handler := market.NewHandler(finder, manager, logger)

		event := queue.EventMarket{
			ID:            "1.234501",
			EventID:       45617,
			Name:          "OVER_UNDER_25",
			CompetitionID: 8,
			SeasonID:      17420,
			EventDate:     time.Unix(1616936636, 0),
			Exchange:      "betfair",
			Runners:       []*queue.Runner{
				{
					ID:         1,
					Name:       "Yes",
					Sort:       1,
					BackPrices: []queue.PriceSize{
						{
							Price: 1.95,
							Size:  500,
						},
					},
					LayPrices:  []queue.PriceSize{
						{
							Price: 1.99,
							Size:  100,
						},
					},
				},
				{
					ID:         2,
					Name:       "No",
					Sort:       1,
					BackPrices: []queue.PriceSize{
						{
							Price: 3.90,
							Size:  500,
						},
					},
					LayPrices:  []queue.PriceSize{
						{
							Price: 3.99,
							Size:  100,
						},
					},
				},
			},
			Timestamp:     1616936636,
		}

		ctx := context.Background()

		stOne := &strategy.Strategy{}

		finder.On("FindMatchingStrategies", ctx, mock.AnythingOfType("*strategy.FinderQuery")).
			Times(4).
			Return(strategyChannel(stOne))

		manager.On("Manage", ctx, mock.AnythingOfType("*trade.Ticket"), mock.AnythingOfType("*strategy.Strategy")).
			Once().
			Return(errors.New("manager error"))

		manager.On("Manage", ctx, mock.AnythingOfType("*trade.Ticket"), mock.AnythingOfType("*strategy.Strategy")).
			Times(3).
			Return(nil)

		handler.HandleEventMarket(ctx, &event)
		
		assert.Equal(t, 1, len(hook.AllEntries()))
	})
}

type MockStrategyFinder struct {
	mock.Mock
}

func (m *MockStrategyFinder) FindMatchingStrategies(ctx context.Context, q *strategy.FinderQuery) <-chan *strategy.Strategy {
	args := m.Called(ctx, q)
	return args.Get(0).(<-chan *strategy.Strategy)
}

type MockTradeManager struct {
	mock.Mock
}

func (m *MockTradeManager) Manage(ctx context.Context, t *trade.Ticket, s *strategy.Strategy) error {
	args := m.Called(ctx, t, s)
	return args.Error(0)
}

func strategyChannel(s *strategy.Strategy) <-chan *strategy.Strategy {
	ch := make(chan *strategy.Strategy, 1)

	ch <- s

	close(ch)

	return ch
}
