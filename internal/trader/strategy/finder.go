package strategy

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type finder struct {
	reader   Reader
	matcher  FilterMatcher
	logger   *logrus.Logger
}

func (h *finder) FindMatchingStrategies(ctx context.Context, q *FinderQuery) <-chan *Strategy {
	ch := make(chan *Strategy, 100)

	go h.findStrategies(ctx, q, ch)

	return ch
}

func (h *finder) findStrategies(ctx context.Context, q *FinderQuery, ch chan<- *Strategy) {
	defer close(ch)

	var wg sync.WaitGroup

	query := ReaderQuery{
		Market:        &q.MarketName,
		Runner:        &q.RunnerName,
		Price:         &q.Price,
		CompetitionID: &q.CompetitionID,
		Side:          &q.Side,
		Status:        &q.Status,
	}

	st, err := h.reader.Get(&query)

	if err != nil {
		h.logger.Errorf("error fetching matches strategies: %+v", err)
		return
	}

	for _, s := range st {
		wg.Add(1)
		h.filterStrategy(ctx, s, q.EventID, ch, &wg)
	}

	wg.Wait()
}

func (h *finder) filterStrategy(ctx context.Context, s *Strategy, eventID uint64, ch chan<- *Strategy, wg *sync.WaitGroup) {
	query := MatcherQuery{
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

func NewFinder(r Reader, f FilterMatcher, l *logrus.Logger) Finder {
	return &finder{
		reader:  r,
		matcher: f,
		logger:  l,
	}
}
