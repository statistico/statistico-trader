package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/classify"
	"sync"
)

type TradeFinder interface {
	Find(ctx context.Context, q *TradeQuery) <-chan *statistico.StrategyTrade
}

type tradeFinder struct {
	matcher classify.FilterMatcher
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
	query := classify.MatcherQuery{
		EventID:       mk.EventId,
		ResultFilters: q.ResultFilters,
		StatFilters:   q.StatFilters,
	}

	matches, err := t.matcher.MatchesFilters(ctx, &query)

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

	if matches {
		ch <- nil
	}

	wg.Done()
}

func NewTradeFinder(m classify.FilterMatcher, l *logrus.Logger) TradeFinder {
	return &tradeFinder{
		matcher: m,
		logger:  l,
	}
}
