package strategy

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-trader/internal/trader/classify"
	"github.com/statistico/statistico-trader/internal/trader/market"
	"sync"
)

type finder struct {
	reader   Reader
	matcher  classify.FilterMatcher
	logger   *logrus.Logger
}

func (h *finder) FindMatchingStrategies(ctx context.Context, m *market.Runner) <-chan *Strategy {
	ch := make(chan *Strategy, 100)

	go h.findStrategies(ctx, m, ch)

	return ch
}

func (h *finder) findStrategies(ctx context.Context, m *market.Runner, ch chan<- *Strategy) {
	defer close(ch)

	var wg sync.WaitGroup

	active := Active

	query := ReaderQuery{
		Market:        &m.MarketName,
		Runner:        &m.RunnerName,
		Price:         &m.Price.Value,
		CompetitionID: &m.CompetitionID,
		Side:          &m.Price.Side,
		Status:        &active,
	}

	st, err := h.reader.Get(&query)

	if err != nil {
		h.logger.Errorf("error fetching matches strategies: %+v", err)
		return
	}

	for _, s := range st {
		wg.Add(1)
		h.filterStrategy(ctx, s, m.EventID, ch, &wg)
	}

	wg.Wait()
}

func (h *finder) filterStrategy(ctx context.Context, s *Strategy, eventID uint64, ch chan<- *Strategy, wg *sync.WaitGroup) {
	query := classify.MatcherQuery{
		EventID:       eventID,
		ResultFilters: nil,
		StatFilters:   nil,
	}

	matches, err := h.matcher.MatchesFilters(ctx, &query)

	if err != nil {
		h.logger.Errorf("error matching strategy %s: %+v", s.ID.String(), err)
		wg.Done()
		return
	}

	if matches {
		ch <- s
	}

	wg.Done()
}

func NewFinder(r Reader, c classify.FilterMatcher, l *logrus.Logger) Finder {
	return &finder{
		reader:  r,
		matcher: c,
		logger:  l,
	}
}
