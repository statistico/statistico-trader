package market

import (
	"context"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data-go-grpc-client"
	"github.com/statistico/statistico-strategy/internal/trader/classify"
	"time"
)

type TradeFactory interface {
	CreateTrade(ctx context.Context, q *Query) (*Trade, error)
}

type tradeFactory struct {
	fixtureClient statisticodata.FixtureClient
	resultClient  statisticodata.ResultClient
	filterMatcher classify.FilterMatcher
	clock         clockwork.Clock
}

func (t *tradeFactory) CreateTrade(ctx context.Context, q *Query) (*Trade, error) {
	fixture, err := t.fixtureClient.ByID(ctx, q.EventId)

	if err != nil {
		return nil, err
	}

	fix := classify.Fixture{
		ID:         uint64(fixture.Id),
		HomeTeamID: fixture.HomeTeam.Id,
		AwayTeamID: fixture.AwayTeam.Id,
		Date:       time.Unix(fixture.DateTime.Utc, 0),
		SeasonID:   fixture.Season.Id,
	}

	matches, err := t.filterMatcher.MatchesFilters(ctx, &fix, q.ResultFilters, q.StatFilters)

	if err != nil {
		return nil, err
	}

	if !matches {
		return nil, nil
	}

	trade := transformQueryToTrade(q)

	if fix.Date.Before(t.clock.Now()) {
		result, err := t.resultClient.ByID(ctx, q.EventId)

		if err != nil {
			return nil, err
		}

		r, err := parseTradeResult(q, result)

		if err != nil {
			return nil, err
		}

		trade.Result = &r
	}

	return trade, nil
}

func NewTradeFactory(f statisticodata.FixtureClient, r statisticodata.ResultClient, fm classify.FilterMatcher, c clockwork.Clock) TradeFactory {
	return &tradeFactory{
		fixtureClient: f,
		resultClient:  r,
		filterMatcher: fm,
		clock:         c,
	}
}
