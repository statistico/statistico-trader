package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader/market"
	"sync"
)

type TradeFinder interface {
	Find(ctx context.Context, q *TradeQuery) <-chan *statistico.StrategyTrade
}

type tradeFinder struct {
	factory market.TradeFactory
	logger  *logrus.Logger
}

func (t *tradeFinder) Find(ctx context.Context, q *TradeQuery) <-chan *statistico.StrategyTrade {
	ch := make(chan *statistico.StrategyTrade, len(q.Markets))

	go t.handleMarkets(ctx, ch, q)

	return ch
}

func (t *tradeFinder) handleMarkets(ctx context.Context, ch chan<- *statistico.StrategyTrade, q *TradeQuery) {
	defer close(ch)

	wg := sync.WaitGroup{}

	for mk := range q.Markets {
		wg.Add(1)
		go t.filterMarket(ctx, ch, mk, q, &wg)
	}

	wg.Wait()
}

func (t *tradeFinder) filterMarket(ctx context.Context, ch chan<- *statistico.StrategyTrade, mk *statistico.MarketRunner, q *TradeQuery, wg *sync.WaitGroup) {
	date, err := ptypes.Timestamp(mk.EventDate)

	if err != nil {
		t.logger.Errorf("Error parsing market event date: Error %s", err.Error())
		wg.Done()
		return
	}

	query := market.Query{
		MarketName:    mk.GetMarketName(),
		RunnerName:    mk.GetRunnerName(),
		RunnerPrice:   mk.GetPrice().GetValue(),
		EventId:       mk.GetEventId(),
		CompetitionId: mk.GetCompetitionId(),
		SeasonId:      mk.GetSeasonId(),
		EventDate:     date,
		Side:          mk.GetPrice().GetSide().String(),
		Exchange:      mk.GetExchange(),
		ResultFilters: q.RunnerFilters,
		StatFilters:   q.StatFilters,
	}

	trade, err := t.factory.CreateTrade(ctx, &query)

	if err != nil {
		t.logger.Errorf(
			"error handling trade for market %s, runner %s and event %d: %s",
			mk.MarketName,
			mk.RunnerName,
			mk.EventId,
			err.Error(),
		)

		wg.Done()
		return
	}

	if trade != nil {
		tr, err := transformTradeResultToStrategyTrade(trade)

		if err != nil {
			t.logger.Errorf(
				"error converting trade for market %s, runner %s and event %d: %s",
				mk.MarketName,
				mk.RunnerName,
				mk.EventId,
				err.Error(),
			)
			return
		}

		ch <- tr
	}

	wg.Done()
}

func NewTradeFinder(f market.TradeFactory, l *logrus.Logger) TradeFinder {
	return &tradeFinder{
		factory: f,
		logger:  l,
	}
}
