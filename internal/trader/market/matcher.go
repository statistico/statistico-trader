package market

import (
	"context"
	"github.com/statistico/statistico-data-go-grpc-client"
	"github.com/statistico/statistico-strategy/internal/trader/classify"
	"time"
)

type Matcher interface {
	Matches(ctx context.Context, q *MatcherQuery) (bool, error)
}

type matcher struct {
	fixtureClient statisticodata.FixtureClient
	filterMatcher classify.FilterMatcher
}

func (m *matcher) Matches(ctx context.Context, q *MatcherQuery) (bool, error) {
	fixture, err := m.fixtureClient.ByID(ctx, q.EventID)

	if err != nil {
		return false, err
	}

	fix := classify.Fixture{
		ID:         uint64(fixture.Id),
		HomeTeamID: fixture.HomeTeam.Id,
		AwayTeamID: fixture.AwayTeam.Id,
		Date:       time.Unix(fixture.DateTime.Utc, 0),
		SeasonID:   fixture.Season.Id,
	}

	return m.filterMatcher.MatchesFilters(ctx, &fix, q.ResultFilters, q.StatFilters)
}

func NewMatcher(f statisticodata.FixtureClient, fm classify.FilterMatcher) Matcher {
	return &matcher{
		fixtureClient: f,
		filterMatcher: fm,
	}
}
