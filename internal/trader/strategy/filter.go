package strategy

import (
	"context"
	"github.com/statistico/statistico-data-go-grpc-client"
	"time"
)

type FilterMatcher interface {
	MatchesFilters(ctx context.Context, q *MatcherQuery) (bool, error)
}

type filterMatcher struct {
	fixtureClient    statisticodata.FixtureClient
	resultClassifier ResultFilterClassifier
	statClassifier   StatFilterClassifier
}

// MatchesFilters receives a MatcherQuery containing trader.ResultFilter and trader.StatFilter slices and determines if
// Fixture matching EventID matches all filters provided.
func (f *filterMatcher) MatchesFilters(ctx context.Context, q *MatcherQuery) (bool, error) {
	fixture, err := f.fixtureClient.ByID(ctx, q.EventID)

	if err != nil {
		return false, err
	}

	fix := Fixture{
		ID:         uint64(fixture.Id),
		HomeTeamID: fixture.HomeTeam.Id,
		AwayTeamID: fixture.AwayTeam.Id,
		Date:       time.Unix(fixture.DateTime.Utc, 0),
		SeasonID:   fixture.Season.Id,
	}

	for _, filter := range q.ResultFilters {
		success, err := f.resultClassifier.MatchesFilter(ctx, &fix, filter)

		if err != nil {
			return false, err
		}

		if !success {
			return false, nil
		}
	}

	for _, filter := range q.StatFilters {
		success, err := f.statClassifier.MatchesFilter(ctx, &fix, filter)

		if err != nil {
			return false, err
		}

		if !success {
			return false, nil
		}
	}

	return true, nil
}

func NewFilterMatcher(f statisticodata.FixtureClient, r ResultFilterClassifier, s StatFilterClassifier) FilterMatcher {
	return &filterMatcher{fixtureClient: f, resultClassifier: r, statClassifier: s}
}
