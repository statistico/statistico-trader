package strategy_test

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestBuilder_Build(t *testing.T) {
	resultFilters := []*strategy.ResultFilter{
		{
			Team:   "HOME_TEAM",
			Result: "WIN",
			Games:  3,
			Venue:  "HOME",
		},
	}

	statFilters := []*strategy.StatFilter{
		{
			Stat:    "GOALS",
			Team:    "HOME_TEAM",
			Action:  "FOR",
			Games:   3,
			Measure: "AVERAGE",
			Metric:  "GTE",
			Value:   4,
			Venue:   "HOME",
		},
	}

	t.Run("parses markets, matches filters and returns a channel of strategy.Trade struct", func(t *testing.T) {
		t.Helper()

		matcher := new(MockFilterMatcher)
		parser := new(MockResultParser)
		marketClient := new(MockMarketClient)
		logger, hook := test.NewNullLogger()

		builder := strategy.NewBuilder(matcher, parser, marketClient, logger)

		ctx := context.Background()

		min := float32(1.50)
		max := float32(2.50)

		query := strategy.BuilderQuery{
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &min,
			MaxOdds:        &max,
			Line:           "CLOSING",
			Side:           "BACK",
			CompetitionIDs: []uint64{8},
			ResultFilters:  resultFilters,
			StatFilters:    statFilters,
		}

		markets := []*statistico.MarketRunner{
			{
				MarketId:             "1.2345",
				MarketName:           "MATCH_ODDS",
				RunnerId:             1,
				RunnerName:           "Home",
				EventId:              1234,
				CompetitionId:        8,
				SeasonId:             17420,
				EventDate:            timestamppb.New(time.Unix(1617126949, 0)),
				Exchange:             "betfair",
				Price:                &statistico.Price{
					Value:                1.95,
					Size:                 500.03,
					Side:                 statistico.SideEnum_BACK,
					Timestamp:            1617126949,
				},
			},
		}

		marketCh := marketChannel(markets)
		errCh := errChan(nil)

		marketReq := mock.MatchedBy(func(r *statistico.MarketRunnerRequest) bool {
			a := assert.New(t)

			a.Equal("MATCH_ODDS", r.GetMarket())
			a.Equal("Home", r.GetRunner())
			a.Equal(min, r.GetMinOdds().GetValue())
			a.Equal(max, r.GetMaxOdds().GetValue())
			a.Equal("CLOSING", r.GetLine())
			a.Equal(statistico.SideEnum_BACK, r.GetSide())
			a.Equal([]uint64{8}, r.GetCompetitionIds())
			a.Nil(r.GetSeasonIds())
			return true
		})

		marketClient.On("MarketRunnerSearch", ctx, marketReq).Return(marketCh, errCh)

		matcherQuery := strategy.MatcherQuery{
			EventID:       1234,
			ResultFilters: resultFilters,
			StatFilters:   statFilters,
		}

		matcher.On("MatchesFilters", ctx, &matcherQuery).Return(true, nil)

		parser.On("Parse", ctx, uint64(1234), "MATCH_ODDS", "Home", "BACK").Return(strategy.Result("SUCCESS"), nil)

		tradeCh := builder.Build(ctx, &query)

		tr := <-tradeCh

		a := assert.New(t)

		a.Equal("MATCH_ODDS", tr.MarketName)
		a.Equal("Home", tr.RunnerName)
		a.Equal(uint64(1234), tr.EventID)
		a.Equal(uint64(8), tr.CompetitionID)
		a.Equal(uint64(17420), tr.SeasonID)
		a.Equal(int64(1617126949), tr.EventDate.Unix())
		a.Equal("betfair", tr.Exchange)
		a.Equal(float32(1.95), tr.Price)
		a.Equal("BACK", tr.Side)
		a.Equal(strategy.Result("SUCCESS"), tr.Result)
		a.Equal(0, len(hook.AllEntries()))

		matcher.AssertExpectations(t)
		marketClient.AssertExpectations(t)
		parser.AssertExpectations(t)
	})

	t.Run("error is logged if error is returned on error channel returned by market client", func(t *testing.T) {
		t.Helper()

		matcher := new(MockFilterMatcher)
		parser := new(MockResultParser)
		marketClient := new(MockMarketClient)
		logger, hook := test.NewNullLogger()

		builder := strategy.NewBuilder(matcher, parser, marketClient, logger)

		ctx := context.Background()

		min := float32(1.50)
		max := float32(2.50)

		query := strategy.BuilderQuery{
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &min,
			MaxOdds:        &max,
			Line:           "CLOSING",
			Side:           "BACK",
			CompetitionIDs: []uint64{8},
			ResultFilters:  resultFilters,
			StatFilters:    statFilters,
		}

		markets := []*statistico.MarketRunner{
			{
				MarketId:             "1.2345",
				MarketName:           "MATCH_ODDS",
				RunnerId:             1,
				RunnerName:           "Home",
				EventId:              1234,
				CompetitionId:        8,
				SeasonId:             17420,
				EventDate:            timestamppb.New(time.Unix(1617126949, 0)),
				Exchange:             "betfair",
				Price:                &statistico.Price{
					Value:                1.95,
					Size:                 500.03,
					Side:                 statistico.SideEnum_BACK,
					Timestamp:            1617126949,
				},
			},
		}

		marketCh := marketChannel(markets)
		errCh := errChan(errors.New("error in market client"))

		marketReq := mock.MatchedBy(func(r *statistico.MarketRunnerRequest) bool {
			a := assert.New(t)

			a.Equal("MATCH_ODDS", r.GetMarket())
			a.Equal("Home", r.GetRunner())
			a.Equal(min, r.GetMinOdds().GetValue())
			a.Equal(max, r.GetMaxOdds().GetValue())
			a.Equal("CLOSING", r.GetLine())
			a.Equal(statistico.SideEnum_BACK, r.GetSide())
			a.Equal([]uint64{8}, r.GetCompetitionIds())
			a.Nil(r.GetSeasonIds())
			return true
		})

		marketClient.On("MarketRunnerSearch", ctx, marketReq).Return(marketCh, errCh)

		matcherQuery := strategy.MatcherQuery{
			EventID:       1234,
			ResultFilters: resultFilters,
			StatFilters:   statFilters,
		}

		matcher.On("MatchesFilters", ctx, &matcherQuery).Return(true, nil)

		parser.On("Parse", ctx, uint64(1234), "MATCH_ODDS", "Home", "BACK").Return(strategy.Result("SUCCESS"), nil)

		tradeCh := builder.Build(ctx, &query)

		tr := <-tradeCh

		a := assert.New(t)

		a.Equal("MATCH_ODDS", tr.MarketName)
		a.Equal("Home", tr.RunnerName)
		a.Equal(uint64(1234), tr.EventID)
		a.Equal(uint64(8), tr.CompetitionID)
		a.Equal(uint64(17420), tr.SeasonID)
		a.Equal(int64(1617126949), tr.EventDate.Unix())
		a.Equal("betfair", tr.Exchange)
		a.Equal(float32(1.95), tr.Price)
		a.Equal("BACK", tr.Side)
		a.Equal(strategy.Result("SUCCESS"), tr.Result)

		a.Equal("error fetching market runners from odds warehouse: error in market client", hook.LastEntry().Message)
		a.Equal(logrus.ErrorLevel, hook.LastEntry().Level)

		matcher.AssertExpectations(t)
		marketClient.AssertExpectations(t)
		parser.AssertExpectations(t)
	})

	t.Run("trade is not pushed into channel if filter matcher returns false", func(t *testing.T) {
		t.Helper()

		matcher := new(MockFilterMatcher)
		parser := new(MockResultParser)
		marketClient := new(MockMarketClient)
		logger, hook := test.NewNullLogger()

		builder := strategy.NewBuilder(matcher, parser, marketClient, logger)

		ctx := context.Background()

		min := float32(1.50)
		max := float32(2.50)

		query := strategy.BuilderQuery{
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &min,
			MaxOdds:        &max,
			Line:           "CLOSING",
			Side:           "BACK",
			CompetitionIDs: []uint64{8},
			ResultFilters:  resultFilters,
			StatFilters:    statFilters,
		}

		markets := []*statistico.MarketRunner{
			{
				MarketId:             "1.2345",
				MarketName:           "MATCH_ODDS",
				RunnerId:             1,
				RunnerName:           "Home",
				EventId:              1234,
				CompetitionId:        8,
				SeasonId:             17420,
				EventDate:            timestamppb.New(time.Unix(1617126949, 0)),
				Exchange:             "betfair",
				Price:                &statistico.Price{
					Value:                1.95,
					Size:                 500.03,
					Side:                 statistico.SideEnum_BACK,
					Timestamp:            1617126949,
				},
			},
		}

		marketCh := marketChannel(markets)
		errCh := errChan(nil)

		marketReq := mock.MatchedBy(func(r *statistico.MarketRunnerRequest) bool {
			a := assert.New(t)

			a.Equal("MATCH_ODDS", r.GetMarket())
			a.Equal("Home", r.GetRunner())
			a.Equal(min, r.GetMinOdds().GetValue())
			a.Equal(max, r.GetMaxOdds().GetValue())
			a.Equal("CLOSING", r.GetLine())
			a.Equal(statistico.SideEnum_BACK, r.GetSide())
			a.Equal([]uint64{8}, r.GetCompetitionIds())
			a.Nil(r.GetSeasonIds())
			return true
		})

		marketClient.On("MarketRunnerSearch", ctx, marketReq).Return(marketCh, errCh)

		matcherQuery := strategy.MatcherQuery{
			EventID:       1234,
			ResultFilters: resultFilters,
			StatFilters:   statFilters,
		}

		matcher.On("MatchesFilters", ctx, &matcherQuery).Return(false, nil)

		parser.AssertNotCalled(t, "Parse")

		tradeCh := builder.Build(ctx, &query)

		tr := <-tradeCh

		assert.Nil(t, tr)
		assert.Equal(t, 0, len(hook.AllEntries()))

		matcher.AssertExpectations(t)
		marketClient.AssertExpectations(t)
	})

	t.Run("error is logged if error is returned by filter matcher", func(t *testing.T) {
		t.Helper()

		matcher := new(MockFilterMatcher)
		parser := new(MockResultParser)
		marketClient := new(MockMarketClient)
		logger, hook := test.NewNullLogger()

		builder := strategy.NewBuilder(matcher, parser, marketClient, logger)

		ctx := context.Background()

		min := float32(1.50)
		max := float32(2.50)

		query := strategy.BuilderQuery{
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &min,
			MaxOdds:        &max,
			Line:           "CLOSING",
			Side:           "BACK",
			CompetitionIDs: []uint64{8},
			ResultFilters:  resultFilters,
			StatFilters:    statFilters,
		}

		markets := []*statistico.MarketRunner{
			{
				MarketId:             "1.2345",
				MarketName:           "MATCH_ODDS",
				RunnerId:             1,
				RunnerName:           "Home",
				EventId:              1234,
				CompetitionId:        8,
				SeasonId:             17420,
				EventDate:            timestamppb.New(time.Unix(1617126949, 0)),
				Exchange:             "betfair",
				Price:                &statistico.Price{
					Value:                1.95,
					Size:                 500.03,
					Side:                 statistico.SideEnum_BACK,
					Timestamp:            1617126949,
				},
			},
		}

		marketCh := marketChannel(markets)
		errCh := errChan(nil)

		marketReq := mock.MatchedBy(func(r *statistico.MarketRunnerRequest) bool {
			a := assert.New(t)

			a.Equal("MATCH_ODDS", r.GetMarket())
			a.Equal("Home", r.GetRunner())
			a.Equal(min, r.GetMinOdds().GetValue())
			a.Equal(max, r.GetMaxOdds().GetValue())
			a.Equal("CLOSING", r.GetLine())
			a.Equal(statistico.SideEnum_BACK, r.GetSide())
			a.Equal([]uint64{8}, r.GetCompetitionIds())
			a.Nil(r.GetSeasonIds())
			return true
		})

		marketClient.On("MarketRunnerSearch", ctx, marketReq).Return(marketCh, errCh)

		matcherQuery := strategy.MatcherQuery{
			EventID:       1234,
			ResultFilters: resultFilters,
			StatFilters:   statFilters,
		}

		matcher.On("MatchesFilters", ctx, &matcherQuery).Return(true, errors.New("error from filter matcher"))

		parser.AssertNotCalled(t, "Parse")

		tradeCh := builder.Build(ctx, &query)

		tr := <-tradeCh

		assert.Nil(t, tr)
		assert.Equal(t, "error handling trade for market MATCH_ODDS, runner Home and event 1234: error from filter matcher", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)

		matcher.AssertExpectations(t)
		marketClient.AssertExpectations(t)
	})

	t.Run("error is logged if error is returned by result parser", func(t *testing.T) {
		t.Helper()

		matcher := new(MockFilterMatcher)
		parser := new(MockResultParser)
		marketClient := new(MockMarketClient)
		logger, hook := test.NewNullLogger()

		builder := strategy.NewBuilder(matcher, parser, marketClient, logger)

		ctx := context.Background()

		min := float32(1.50)
		max := float32(2.50)

		query := strategy.BuilderQuery{
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &min,
			MaxOdds:        &max,
			Line:           "CLOSING",
			Side:           "BACK",
			CompetitionIDs: []uint64{8},
			ResultFilters:  resultFilters,
			StatFilters:    statFilters,
		}

		markets := []*statistico.MarketRunner{
			{
				MarketId:             "1.2345",
				MarketName:           "MATCH_ODDS",
				RunnerId:             1,
				RunnerName:           "Home",
				EventId:              1234,
				CompetitionId:        8,
				SeasonId:             17420,
				EventDate:            timestamppb.New(time.Unix(1617126949, 0)),
				Exchange:             "betfair",
				Price:                &statistico.Price{
					Value:                1.95,
					Size:                 500.03,
					Side:                 statistico.SideEnum_BACK,
					Timestamp:            1617126949,
				},
			},
		}

		marketCh := marketChannel(markets)
		errCh := errChan(nil)

		marketReq := mock.MatchedBy(func(r *statistico.MarketRunnerRequest) bool {
			a := assert.New(t)

			a.Equal("MATCH_ODDS", r.GetMarket())
			a.Equal("Home", r.GetRunner())
			a.Equal(min, r.GetMinOdds().GetValue())
			a.Equal(max, r.GetMaxOdds().GetValue())
			a.Equal("CLOSING", r.GetLine())
			a.Equal(statistico.SideEnum_BACK, r.GetSide())
			a.Equal([]uint64{8}, r.GetCompetitionIds())
			a.Nil(r.GetSeasonIds())
			return true
		})

		marketClient.On("MarketRunnerSearch", ctx, marketReq).Return(marketCh, errCh)

		matcherQuery := strategy.MatcherQuery{
			EventID:       1234,
			ResultFilters: resultFilters,
			StatFilters:   statFilters,
		}

		matcher.On("MatchesFilters", ctx, &matcherQuery).Return(true, nil)

		parser.On("Parse", ctx, uint64(1234), "MATCH_ODDS", "Home", "BACK").Return(strategy.Result("SUCCESS"), errors.New("parser error"))

		tradeCh := builder.Build(ctx, &query)

		tr := <-tradeCh

		assert.Nil(t, tr)
		assert.Equal(t, "error handling trade for market MATCH_ODDS, runner Home and event 1234: parser error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)

		matcher.AssertExpectations(t)
		marketClient.AssertExpectations(t)
		parser.AssertExpectations(t)
	})
}

type MockResultParser struct {
	mock.Mock
}

func (m *MockResultParser) Parse(ctx context.Context, eventID uint64, market, runner, side string) (strategy.Result, error) {
	args := m.Called(ctx, eventID, market, runner, side)
	return args.Get(0).(strategy.Result), args.Error(1)
}

type MockMarketClient struct {
	mock.Mock
}

func (m *MockMarketClient) MarketRunnerSearch(ctx context.Context, r *statistico.MarketRunnerRequest) (<-chan *statistico.MarketRunner, <-chan error) {
	args := m.Called(ctx, r)
	return args.Get(0).(<-chan *statistico.MarketRunner), args.Get(1).(<-chan error)
}

func marketChannel(markets []*statistico.MarketRunner) <-chan *statistico.MarketRunner {
	ch := make(chan *statistico.MarketRunner, len(markets))

	for _, m := range markets {
		ch <- m
	}

	close(ch)

	return ch
}

func errChan(e error) <-chan error {
	ch := make(chan error, 1)
	ch <- e
	close(ch)
	return ch
}
