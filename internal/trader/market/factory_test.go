package market_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/classify"
	"github.com/statistico/statistico-strategy/internal/trader/market"
	m "github.com/statistico/statistico-strategy/internal/trader/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTradeFactory_CreateTrade(t *testing.T) {
	ctx := context.Background()

	query := market.Query{
		MarketName:    "MATCH_ODDS",
		RunnerName:    "Home",
		RunnerPrice:   1.65,
		EventId:       192810,
		CompetitionId: 8,
		SeasonId:      17420,
		EventDate:     time.Unix(1616052304, 0),
		Side:          "BACK",
		Exchange:      "betfair",
		ResultFilters: nil,
		StatFilters:   nil,
	}

	fixture := statistico.Fixture{
		Id:          192810,
		Competition: &statistico.Competition{Id: 8},
		Season:      &statistico.Season{Id: 17420},
		HomeTeam:    &statistico.Team{Id: 5},
		AwayTeam:    &statistico.Team{Id: 10},
		DateTime:    &statistico.Date{Utc: 1616052304},
	}

	t.Run("returns a Trade struct", func(t *testing.T) {
		t.Helper()

		fixtureClient := new(m.FixtureClient)
		resultClient := new(m.ResultClient)
		matcher := new(MockFilterMatcher)
		clock := clockwork.NewFakeClockAt(time.Unix(1616034304, 0))

		factory := market.NewTradeFactory(fixtureClient, resultClient, matcher, clock)

		fixtureClient.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		cf := classify.Fixture{
			ID:         192810,
			HomeTeamID: 5,
			AwayTeamID: 10,
			Date:       time.Unix(1616052304, 0),
			SeasonID:   17420,
		}

		matcher.On("MatchesFilters", ctx, &cf, query.ResultFilters, query.StatFilters).Return(true, nil)

		trade, err := factory.CreateTrade(ctx, &query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		resultClient.AssertNotCalled(t, "ByID")

		a := assert.New(t)
		a.Equal("MATCH_ODDS", trade.MarketName)
		a.Equal("Home", trade.RunnerName)
		a.Equal(float32(1.65), trade.RunnerPrice)
		a.Equal(uint64(192810), trade.EventId)
		a.Equal(uint64(8), trade.CompetitionId)
		a.Equal(uint64(17420), trade.SeasonId)
		a.Equal(time.Unix(1616052304, 0), trade.EventDate)
		a.Equal("BACK", trade.Side)
		a.Nil(trade.Result)
	})

	t.Run("returns a Trade struct containing result if fixture data is after current date", func(t *testing.T) {
		t.Helper()

		fixtureClient := new(m.FixtureClient)
		resultClient := new(m.ResultClient)
		matcher := new(MockFilterMatcher)
		clock := clockwork.NewFakeClockAt(time.Unix(1616054304, 0))

		factory := market.NewTradeFactory(fixtureClient, resultClient, matcher, clock)

		fixtureClient.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		cf := classify.Fixture{
			ID:         192810,
			HomeTeamID: 5,
			AwayTeamID: 10,
			Date:       time.Unix(1616052304, 0),
			SeasonID:   17420,
		}

		matcher.On("MatchesFilters", ctx, &cf, query.ResultFilters, query.StatFilters).Return(true, nil)

		result := statistico.Result{
			Id:       192810,
			HomeTeam: &statistico.Team{Id: 5},
			AwayTeam: &statistico.Team{Id: 10},
			Stats: &statistico.MatchStats{
				HomeScore: &wrappers.UInt32Value{Value: 2},
				AwayScore: &wrappers.UInt32Value{Value: 1},
			},
		}

		resultClient.On("ByID", ctx, uint64(192810)).Return(&result, nil)

		trade, err := factory.CreateTrade(ctx, &query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal("MATCH_ODDS", trade.MarketName)
		a.Equal("Home", trade.RunnerName)
		a.Equal(float32(1.65), trade.RunnerPrice)
		a.Equal(uint64(192810), trade.EventId)
		a.Equal(uint64(8), trade.CompetitionId)
		a.Equal(uint64(17420), trade.SeasonId)
		a.Equal(time.Unix(1616052304, 0), trade.EventDate)
		a.Equal("BACK", trade.Side)
		a.Equal("SUCCESS", *trade.Result)
	})

	t.Run("returns an error if error is returned by fixture client", func(t *testing.T) {
		t.Helper()

		fixtureClient := new(m.FixtureClient)
		resultClient := new(m.ResultClient)
		matcher := new(MockFilterMatcher)
		clock := clockwork.NewFakeClockAt(time.Unix(1616054304, 0))

		factory := market.NewTradeFactory(fixtureClient, resultClient, matcher, clock)

		fixtureClient.On("ByID", ctx, uint64(192810)).Return(&statistico.Fixture{}, errors.New("error occurred"))

		matcher.AssertNotCalled(t, "MatchesFilters")
		resultClient.AssertNotCalled(t, "ByID")

		_, err := factory.CreateTrade(ctx, &query)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("returns error if error returns by FilterMatcher", func(t *testing.T) {
		t.Helper()

		fixtureClient := new(m.FixtureClient)
		resultClient := new(m.ResultClient)
		matcher := new(MockFilterMatcher)
		clock := clockwork.NewFakeClockAt(time.Unix(1616054304, 0))

		factory := market.NewTradeFactory(fixtureClient, resultClient, matcher, clock)

		fixtureClient.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		cf := classify.Fixture{
			ID:         192810,
			HomeTeamID: 5,
			AwayTeamID: 10,
			Date:       time.Unix(1616052304, 0),
			SeasonID:   17420,
		}

		matcher.On("MatchesFilters", ctx, &cf, query.ResultFilters, query.StatFilters).Return(false, errors.New("error occurred"))

		resultClient.AssertNotCalled(t, "ByID")

		_, err := factory.CreateTrade(ctx, &query)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("returns an error if error returned by result client", func(t *testing.T) {
		t.Helper()

		fixtureClient := new(m.FixtureClient)
		resultClient := new(m.ResultClient)
		matcher := new(MockFilterMatcher)
		clock := clockwork.NewFakeClockAt(time.Unix(1616054304, 0))

		factory := market.NewTradeFactory(fixtureClient, resultClient, matcher, clock)

		fixtureClient.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		cf := classify.Fixture{
			ID:         192810,
			HomeTeamID: 5,
			AwayTeamID: 10,
			Date:       time.Unix(1616052304, 0),
			SeasonID:   17420,
		}

		matcher.On("MatchesFilters", ctx, &cf, query.ResultFilters, query.StatFilters).Return(true, nil)

		resultClient.On("ByID", ctx, uint64(192810)).Return(&statistico.Result{}, errors.New("error occurred"))

		_, err := factory.CreateTrade(ctx, &query)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

type MockFilterMatcher struct {
	mock.Mock
}

func (m *MockFilterMatcher) MatchesFilters(ctx context.Context, fix *classify.Fixture, r []*trader.ResultFilter, s []*trader.StatFilter) (bool, error) {
	args := m.Called(ctx, fix, r, s)
	return args.Get(0).(bool), args.Error(1)
}
