package strategy_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFinder_FindMatchingStrategies(t *testing.T) {
	t.Run("returns a channel of matching strategy.Strategy struct", func(t *testing.T) {
		t.Helper()

		reader := new(MockStrategyReader)
		matcher := new(MockFilterMatcher)
		logger, _ := test.NewNullLogger()

		finder := strategy.NewFinder(reader, matcher, logger)

		ctx := context.Background()

		query := strategy.FinderQuery{
			MarketName:    "MATCH_ODDS",
			RunnerName:    "Home",
			EventID:       1234,
			CompetitionID: 8,
			Price:         1.95,
			Side:          "BACK",
			Status:        "ACTIVE",
		}

		mockReaderQuery := mock.MatchedBy(func(q *strategy.ReaderQuery) bool {
			a := assert.New(t)

			a.Equal(query.MarketName, *q.Market)
			a.Equal(query.RunnerName, *q.Runner)
			a.Equal(query.Price, *q.Price)
			a.Equal(query.CompetitionID, *q.CompetitionID)
			a.Equal(query.Side, *q.Side)
			a.Equal(query.Status, *q.Status)
			return true
		})

		stOne := &strategy.Strategy{
			ID: uuid.New(),
			Name: "Strategy One",
		}

		strategies := []*strategy.Strategy{stOne}

		reader.On("Get", mockReaderQuery).Return(strategies, nil)

		matcherQuery := mock.MatchedBy(func(q *strategy.MatcherQuery) bool {
			a := assert.New(t)

			a.Equal(uint64(1234), q.EventID)
			a.Equal(stOne.ResultFilters, q.ResultFilters)
			a.Equal(stOne.StatFilters, q.StatFilters)
			return true
		})

		matcher.On("MatchesFilters", ctx, matcherQuery).Return(true, nil)

		ch := finder.FindMatchingStrategies(ctx, &query)

		fetched := <- ch

		assert.Equal(t, stOne, fetched)
	})

	t.Run("does not push strategy into channel if matcher returns false", func(t *testing.T) {
		t.Helper()

		reader := new(MockStrategyReader)
		matcher := new(MockFilterMatcher)
		logger, _ := test.NewNullLogger()

		finder := strategy.NewFinder(reader, matcher, logger)

		ctx := context.Background()

		query := strategy.FinderQuery{
			MarketName:    "MATCH_ODDS",
			RunnerName:    "Home",
			EventID:       1234,
			CompetitionID: 8,
			Price:         1.95,
			Side:          "BACK",
			Status:        "ACTIVE",
		}

		mockReaderQuery := mock.MatchedBy(func(q *strategy.ReaderQuery) bool {
			a := assert.New(t)

			a.Equal(query.MarketName, *q.Market)
			a.Equal(query.RunnerName, *q.Runner)
			a.Equal(query.Price, *q.Price)
			a.Equal(query.CompetitionID, *q.CompetitionID)
			a.Equal(query.Side, *q.Side)
			a.Equal(query.Status, *q.Status)
			return true
		})

		stOne := &strategy.Strategy{
			ID: uuid.New(),
			Name: "Strategy One",
		}

		strategies := []*strategy.Strategy{stOne}

		reader.On("Get", mockReaderQuery).Return(strategies, nil)

		matcherQuery := mock.MatchedBy(func(q *strategy.MatcherQuery) bool {
			a := assert.New(t)

			a.Equal(uint64(1234), q.EventID)
			a.Equal(stOne.ResultFilters, q.ResultFilters)
			a.Equal(stOne.StatFilters, q.StatFilters)
			return true
		})

		matcher.On("MatchesFilters", ctx, matcherQuery).Return(true, nil)

		ch := finder.FindMatchingStrategies(ctx, &query)

		assert.Equal(t, 0, len(ch))
	})

	t.Run("error is logged if returned by strategy.Reader", func(t *testing.T) {
		t.Helper()

		reader := new(MockStrategyReader)
		matcher := new(MockFilterMatcher)
		logger, hook := test.NewNullLogger()

		finder := strategy.NewFinder(reader, matcher, logger)

		ctx := context.Background()

		query := strategy.FinderQuery{
			MarketName:    "MATCH_ODDS",
			RunnerName:    "Home",
			EventID:       1234,
			CompetitionID: 8,
			Price:         1.95,
			Side:          "BACK",
			Status:        "ACTIVE",
		}

		mockReaderQuery := mock.MatchedBy(func(q *strategy.ReaderQuery) bool {
			a := assert.New(t)

			a.Equal(query.MarketName, *q.Market)
			a.Equal(query.RunnerName, *q.Runner)
			a.Equal(query.Price, *q.Price)
			a.Equal(query.CompetitionID, *q.CompetitionID)
			a.Equal(query.Side, *q.Side)
			a.Equal(query.Status, *q.Status)
			return true
		})

		reader.On("Get", mockReaderQuery).Return([]*strategy.Strategy{}, errors.New("reader error"))

		matcher.AssertNotCalled(t, "MatchesFilters")

		ch := finder.FindMatchingStrategies(ctx, &query)

		<-ch

		assert.Equal(t, "error fetching matches strategies: reader error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("error is logged if returned by FilterMatcher", func(t *testing.T) {
		t.Helper()

		reader := new(MockStrategyReader)
		matcher := new(MockFilterMatcher)
		logger, hook := test.NewNullLogger()

		finder := strategy.NewFinder(reader, matcher, logger)

		ctx := context.Background()

		query := strategy.FinderQuery{
			MarketName:    "MATCH_ODDS",
			RunnerName:    "Home",
			EventID:       1234,
			CompetitionID: 8,
			Price:         1.95,
			Side:          "BACK",
			Status:        "ACTIVE",
		}

		mockReaderQuery := mock.MatchedBy(func(q *strategy.ReaderQuery) bool {
			a := assert.New(t)

			a.Equal(query.MarketName, *q.Market)
			a.Equal(query.RunnerName, *q.Runner)
			a.Equal(query.Price, *q.Price)
			a.Equal(query.CompetitionID, *q.CompetitionID)
			a.Equal(query.Side, *q.Side)
			a.Equal(query.Status, *q.Status)
			return true
		})

		stOne := &strategy.Strategy{
			ID: uuid.MustParse("c1c53e13-bded-46d5-8fe5-01088262efb5"),
			Name: "Strategy One",
		}

		strategies := []*strategy.Strategy{stOne}

		reader.On("Get", mockReaderQuery).Return(strategies, nil)

		matcherQuery := mock.MatchedBy(func(q *strategy.MatcherQuery) bool {
			a := assert.New(t)

			a.Equal(uint64(1234), q.EventID)
			a.Equal(stOne.ResultFilters, q.ResultFilters)
			a.Equal(stOne.StatFilters, q.StatFilters)
			return true
		})

		matcher.On("MatchesFilters", ctx, matcherQuery).Return(false, errors.New("matcher error"))

		ch := finder.FindMatchingStrategies(ctx, &query)

		<-ch

		assert.Equal(t, "error matching strategy c1c53e13-bded-46d5-8fe5-01088262efb5: matcher error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})
}

type MockStrategyReader struct {
	mock.Mock
}

func (m *MockStrategyReader) Get(q *strategy.ReaderQuery) ([]*strategy.Strategy, error) {
	args := m.Called(q)
	return args.Get(0).([]*strategy.Strategy), args.Error(1)
}

type MockFilterMatcher struct {
	mock.Mock
}

func (m *MockFilterMatcher) MatchesFilters(ctx context.Context, q *strategy.MatcherQuery) (bool, error) {
	args := m.Called(ctx, q)
	return args.Get(0).(bool), args.Error(1)
}
