package classify

import (
	"context"
	"github.com/statistico/statistico-strategy/internal/trader"
)

type FilterMatcher interface {
	MatchesFilters(ctx context.Context, fix *Fixture, r []*trader.ResultFilter, s []*trader.StatFilter) (bool, error)
}

type filterMatcher struct {
	resultClassifier ResultFilterClassifier
	statClassifier   StatFilterClassifier
}

// MatchesFilters receives a Fixture along with trader.ResultFilter and trader.StatFilter slices and determines if
// Fixture matches all filters provided.
func (f *filterMatcher) MatchesFilters(ctx context.Context, fix *Fixture, r []*trader.ResultFilter, s []*trader.StatFilter) (bool, error) {
	for _, filter := range r {
		success, err := f.resultClassifier.MatchesFilter(ctx, fix, filter)

		if err != nil {
			return false, err
		}

		if !success {
			return false, nil
		}
	}

	for _, filter := range s {
		success, err := f.statClassifier.MatchesFilter(ctx, fix, filter)

		if err != nil {
			return false, err
		}

		if !success {
			return false, nil
		}
	}

	return true, nil
}

func NewFilterMatcher(r ResultFilterClassifier, s StatFilterClassifier) FilterMatcher {
	return &filterMatcher{resultClassifier: r, statClassifier: s}
}
