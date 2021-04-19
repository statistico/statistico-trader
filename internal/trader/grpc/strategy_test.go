package grpc_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-proto/go"
	errors2 "github.com/statistico/statistico-trader/internal/trader/errors"
	g "github.com/statistico/statistico-trader/internal/trader/grpc"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestStrategyService_BuildStrategy(t *testing.T) {
	req := statistico.BuildStrategyRequest{
		Market:         "BOTH_TEAMS_TO_SCORE",
		Runner:         "Yes",
		Line:           "CLOSING",
		Side:           statistico.SideEnum_BACK,
		MinOdds:        &wrappers.FloatValue{Value: 1.95},
		MaxOdds:        &wrappers.FloatValue{Value: 3.55},
		CompetitionIds: []uint64{8, 9, 10},
		SeasonIds:      []uint64{24, 25, 26},
		DateFrom:       &timestamp.Timestamp{Seconds: 1584014400},
		DateTo:         &timestamp.Timestamp{Seconds: 1584014400},
		ResultFilters: []*statistico.ResultFilter{
			{
				Team:   statistico.TeamEnum_HOME_TEAM,
				Result: statistico.ResultEnum_WIN_DRAW,
				Games:  2,
				Venue:  statistico.VenueEnum_HOME_AWAY,
			},
		},
		StatFilters: []*statistico.StatFilter{
			{
				Stat:    statistico.StatEnum_GOALS,
				Team:    statistico.TeamEnum_HOME_TEAM,
				Action:  statistico.ActionEnum_AGAINST,
				Games:   4,
				Metric:  statistico.MetricEnum_GTE,
				Measure: statistico.MeasureEnum_AVERAGE,
				Value:   3.1,
				Venue:   statistico.VenueEnum_AWAY,
			},
		},
	}

	query := mock.MatchedBy(func(q *strategy.BuilderQuery) bool {
		resFil := []*strategy.ResultFilter{
			{
				Team:   "HOME_TEAM",
				Result: "WIN_DRAW",
				Games:  uint8(2),
				Venue:  "HOME_AWAY",
			},
		}

		statFil := []*strategy.StatFilter{
			{
				Stat:    "GOALS",
				Team:    "HOME_TEAM",
				Action:  "AGAINST",
				Games:   uint8(4),
				Metric:  "GTE",
				Measure: "AVERAGE",
				Value:   3.1,
				Venue:   "AWAY",
			},
		}

		a := assert.New(t)

		a.Equal("BOTH_TEAMS_TO_SCORE", q.Market)
		a.Equal("Yes", q.Runner)
		a.Equal("CLOSING", q.Line)
		a.Equal("BACK", q.Side)
		a.Equal([]uint64{8, 9, 10}, q.CompetitionIDs)
		a.Equal([]uint64{24, 25, 26}, q.SeasonIDs)
		a.Equal(float32(1.95), *q.MinOdds)
		a.Equal(float32(3.55), *q.MaxOdds)
		a.Equal(resFil, q.ResultFilters)
		a.Equal(statFil, q.StatFilters)

		return true
	})

	trades := []*strategy.Trade{
		{
			MarketName:    "BOTH_TEAMS_TO_SCORE",
			RunnerName:    "Yes",
			Price:   	   1.95,
			EventID:       138171,
			CompetitionID: 8,
			SeasonID:      17420,
			Exchange: 	   "betfair",
			Side:           "BACK",
			EventDate:     time.Unix(1584014400, 0),
			Result:        strategy.Result("SUCCESS"),
		},
	}

	t.Run("builds strategy using builder and stream statistico.StrategyTrade structs", func(t *testing.T) {
		t.Helper()

		writer := new(MockStrategyWriter)
		reader := new(MockStrategyReader)
		builder := new(MockStrategyBuilder)
		logger, hook := test.NewNullLogger()
		clock := clockwork.NewFakeClockAt(time.Unix(1616936636, 0))

		stream := new(MockStrategyBuildServer)

		service := g.NewStrategyService(builder, writer, reader, logger, clock)

		ctx := context.Background()

		stream.On("Context").Return(ctx)

		tradeCh := tradeChannel(trades)

		builder.On("Build", ctx, query).Return(tradeCh)

		stream.On("Send", mock.AnythingOfType("*statistico.StrategyTrade")).Once().Return(nil)

		err := service.BuildStrategy(&req, stream)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 0, len(hook.Entries))
		builder.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("logs error if error streaming StrategyResult struct", func(t *testing.T) {
		t.Helper()

		writer := new(MockStrategyWriter)
		reader := new(MockStrategyReader)
		builder := new(MockStrategyBuilder)
		logger, hook := test.NewNullLogger()
		clock := clockwork.NewFakeClockAt(time.Unix(1616936636, 0))

		stream := new(MockStrategyBuildServer)

		service := g.NewStrategyService(builder, writer, reader, logger, clock)

		ctx := context.Background()

		stream.On("Context").Return(ctx)

		tradeCh := tradeChannel(trades)

		builder.On("Build", ctx, query).Return(tradeCh)

		stream.On("Send", mock.AnythingOfType("*statistico.StrategyTrade")).Once().Return(errors.New("stream error"))

		err := service.BuildStrategy(&req, stream)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, "error streaming strategy trade back to client: stream error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		stream.AssertExpectations(t)
	})
}

func TestStrategyService_SaveStrategy(t *testing.T) {
	t.Run("saves and returns a new strategy", func(t *testing.T) {
		t.Helper()

		writer := new(MockStrategyWriter)
		reader := new(MockStrategyReader)
		builder := new(MockStrategyBuilder)
		logger, _ := test.NewNullLogger()
		clock := clockwork.NewFakeClockAt(time.Unix(1616936636, 0))

		service := g.NewStrategyService(builder, writer, reader, logger, clock)

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &wrappers.FloatValue{Value: 1.50},
			MaxOdds:        &wrappers.FloatValue{Value: 5.25},
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			ResultFilters: []*statistico.ResultFilter{
				{
					Team:   statistico.TeamEnum_HOME_TEAM,
					Result: statistico.ResultEnum_WIN_DRAW,
					Games:  2,
					Venue:  statistico.VenueEnum_HOME_AWAY,
				},
			},
			StatFilters: []*statistico.StatFilter{
				{
					Stat:    statistico.StatEnum_GOALS,
					Team:    statistico.TeamEnum_HOME_TEAM,
					Action:  statistico.ActionEnum_AGAINST,
					Games:   4,
					Metric:  statistico.MetricEnum_GTE,
					Measure: statistico.MeasureEnum_AVERAGE,
					Value:   3.1,
					Venue:   statistico.VenueEnum_AWAY,
				},
			},
			Visibility: statistico.VisibilityEnum_PRIVATE,
			StakingPlan: &statistico.StakingPlan{
				Name:  statistico.StakingPlanEnum_PERCENTAGE,
				Value: 2.5,
			},
		}

		st := mock.MatchedBy(func(s *strategy.Strategy) bool {
			res := []*strategy.ResultFilter{
				{
					Team:   "HOME_TEAM",
					Result: "WIN_DRAW",
					Games:  2,
					Venue:  "HOME_AWAY",
				},
			}

			stat := []*strategy.StatFilter{
				{
					Stat:    "GOALS",
					Team:    "HOME_TEAM",
					Action:  "AGAINST",
					Games:   4,
					Measure: "AVERAGE",
					Metric:  "GTE",
					Value:   3.1,
					Venue:   "AWAY",
				},
			}

			a := assert.New(t)

			a.Equal(r.GetName(), s.Name)
			a.Equal(r.GetDescription(), s.Description)
			a.Equal(r.GetMarket(), s.MarketName)
			a.Equal(r.GetRunner(), s.RunnerName)
			a.Equal(r.GetCompetitionIds(), s.CompetitionIDs)
			a.Equal(r.GetSide().String(), s.Side)
			a.Equal(r.GetVisibility().String(), s.Visibility)
			a.Equal("ACTIVE", s.Status)
			a.Equal(res, s.ResultFilters)
			a.Equal(stat, s.StatFilters)
			a.Equal(time.Unix(1616936636, 0), s.CreatedAt)
			a.Equal(time.Unix(1616936636, 0), s.UpdatedAt)
			return true
		})

		writer.On("Insert", st).Return(nil)

		ctx := context.WithValue(context.Background(), "userID", "a5f04fd2-dfe7-41c1-af38-d490119705d8")

		s, err := service.SaveStrategy(ctx, r)

		if err != nil {
			t.Fatalf("Expected error, got %s", err.Error())
		}

		a := assert.New(t)

		date := timestamppb.New(time.Unix(1616936636, 0))

		a.Equal(r.GetName(), s.GetName())
		a.Equal(r.GetDescription(), s.GetDescription())
		a.Equal("a5f04fd2-dfe7-41c1-af38-d490119705d8", s.GetUserId())
		a.Equal(r.GetMarket(), s.GetMarket())
		a.Equal(r.GetRunner(), s.GetRunner())
		a.Equal(r.GetMinOdds(), s.GetMinOdds())
		a.Equal(r.GetMaxOdds(), s.GetMaxOdds())
		a.Equal(r.GetCompetitionIds(), s.GetCompetitionIds())
		a.Equal(r.GetSide(), s.GetSide())
		a.Equal(r.GetVisibility(), s.GetVisibility())
		a.Equal(statistico.StrategyStatusEnum_ACTIVE, s.GetStatus())
		a.Equal(r.GetResultFilters(), s.GetResultFilters())
		a.Equal(r.GetStatFilters(), s.GetStatFilters())
		a.Equal(date, s.GetCreatedAt())
		a.Equal(date, s.GetUpdatedAt())

		writer.AssertExpectations(t)
	})

	t.Run("returns an invalid argument error if User ID provided is not a valid uuid string", func(t *testing.T) {
		t.Helper()

		writer := new(MockStrategyWriter)
		reader := new(MockStrategyReader)
		builder := new(MockStrategyBuilder)
		logger, _ := test.NewNullLogger()
		clock := clockwork.NewFakeClockAt(time.Unix(1616936636, 0))

		service := g.NewStrategyService(builder, writer, reader, logger, clock)

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &wrappers.FloatValue{Value: 1.50},
			MaxOdds:        &wrappers.FloatValue{Value: 5.25},
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			ResultFilters: []*statistico.ResultFilter{
				{
					Team:   statistico.TeamEnum_HOME_TEAM,
					Result: statistico.ResultEnum_WIN_DRAW,
					Games:  2,
					Venue:  statistico.VenueEnum_HOME_AWAY,
				},
			},
			StatFilters: []*statistico.StatFilter{
				{
					Stat:    statistico.StatEnum_GOALS,
					Team:    statistico.TeamEnum_HOME_TEAM,
					Action:  statistico.ActionEnum_AGAINST,
					Games:   4,
					Metric:  statistico.MetricEnum_GTE,
					Measure: statistico.MeasureEnum_AVERAGE,
					Value:   3.1,
					Venue:   statistico.VenueEnum_AWAY,
				},
			},
			Visibility: statistico.VisibilityEnum_PRIVATE,
			StakingPlan: &statistico.StakingPlan{
				Name:  statistico.StakingPlanEnum_PERCENTAGE,
				Value: 2.5,
			},
		}

		ctx := context.WithValue(context.Background(), "userID", "a")

		_, err := service.SaveStrategy(ctx, r)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if e, ok := status.FromError(err); ok {
			assert.Equal(t, codes.Code(3), e.Code())
		}

		assert.Equal(t, "rpc error: code = InvalidArgument desc = user id provided is not a valid uuid string: invalid UUID length: 1", err.Error())

		writer.AssertNotCalled(t, "Insert")
	})

	t.Run("returns duplication error if DuplicationError returned by StrategyWriter", func(t *testing.T) {
		t.Helper()

		writer := new(MockStrategyWriter)
		reader := new(MockStrategyReader)
		builder := new(MockStrategyBuilder)
		logger, _ := test.NewNullLogger()
		clock := clockwork.NewFakeClockAt(time.Unix(1616936636, 0))

		service := g.NewStrategyService(builder, writer, reader, logger, clock)

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &wrappers.FloatValue{Value: 1.50},
			MaxOdds:        &wrappers.FloatValue{Value: 5.25},
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			ResultFilters: []*statistico.ResultFilter{
				{
					Team:   statistico.TeamEnum_HOME_TEAM,
					Result: statistico.ResultEnum_WIN_DRAW,
					Games:  2,
					Venue:  statistico.VenueEnum_HOME_AWAY,
				},
			},
			StatFilters: []*statistico.StatFilter{
				{
					Stat:    statistico.StatEnum_GOALS,
					Team:    statistico.TeamEnum_HOME_TEAM,
					Action:  statistico.ActionEnum_AGAINST,
					Games:   4,
					Metric:  statistico.MetricEnum_GTE,
					Measure: statistico.MeasureEnum_AVERAGE,
					Value:   3.1,
					Venue:   statistico.VenueEnum_AWAY,
				},
			},
			Visibility: statistico.VisibilityEnum_PRIVATE,
			StakingPlan: &statistico.StakingPlan{
				Name:  statistico.StakingPlanEnum_PERCENTAGE,
				Value: 2.5,
			},
		}

		st := mock.MatchedBy(func(s *strategy.Strategy) bool {
			res := []*strategy.ResultFilter{
				{
					Team:   "HOME_TEAM",
					Result: "WIN_DRAW",
					Games:  2,
					Venue:  "HOME_AWAY",
				},
			}

			stat := []*strategy.StatFilter{
				{
					Stat:    "GOALS",
					Team:    "HOME_TEAM",
					Action:  "AGAINST",
					Games:   4,
					Measure: "AVERAGE",
					Metric:  "GTE",
					Value:   3.1,
					Venue:   "AWAY",
				},
			}

			a := assert.New(t)

			a.Equal(r.GetName(), s.Name)
			a.Equal(r.GetDescription(), s.Description)
			a.Equal(r.GetMarket(), s.MarketName)
			a.Equal(r.GetRunner(), s.RunnerName)
			a.Equal(r.GetCompetitionIds(), s.CompetitionIDs)
			a.Equal(r.GetSide().String(), s.Side)
			a.Equal(r.GetVisibility().String(), s.Visibility)
			a.Equal("ACTIVE", s.Status)
			a.Equal(res, s.ResultFilters)
			a.Equal(stat, s.StatFilters)
			a.Equal(time.Unix(1616936636, 0), s.CreatedAt)
			a.Equal(time.Unix(1616936636, 0), s.UpdatedAt)
			return true
		})

		e := errors2.DuplicationError{Message: "Strategy exists with name"}

		writer.On("Insert", st).Return(&e)

		ctx := context.WithValue(context.Background(), "userID", "a5f04fd2-dfe7-41c1-af38-d490119705d8")

		_, err := service.SaveStrategy(ctx, r)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if e, ok := status.FromError(err); ok {
			assert.Equal(t, codes.Code(6), e.Code())
		}

		assert.Equal(t, "rpc error: code = AlreadyExists desc = Duplication error: Strategy exists with name", err.Error())

		writer.AssertExpectations(t)
	})

	t.Run("returns internal server error if non DuplicationError returned by StrategyWriter", func(t *testing.T) {
		t.Helper()

		writer := new(MockStrategyWriter)
		reader := new(MockStrategyReader)
		builder := new(MockStrategyBuilder)
		logger, _ := test.NewNullLogger()
		clock := clockwork.NewFakeClockAt(time.Unix(1616936636, 0))

		service := g.NewStrategyService(builder, writer, reader, logger, clock)

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &wrappers.FloatValue{Value: 1.50},
			MaxOdds:        &wrappers.FloatValue{Value: 5.25},
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			ResultFilters: []*statistico.ResultFilter{
				{
					Team:   statistico.TeamEnum_HOME_TEAM,
					Result: statistico.ResultEnum_WIN_DRAW,
					Games:  2,
					Venue:  statistico.VenueEnum_HOME_AWAY,
				},
			},
			StatFilters: []*statistico.StatFilter{
				{
					Stat:    statistico.StatEnum_GOALS,
					Team:    statistico.TeamEnum_HOME_TEAM,
					Action:  statistico.ActionEnum_AGAINST,
					Games:   4,
					Metric:  statistico.MetricEnum_GTE,
					Measure: statistico.MeasureEnum_AVERAGE,
					Value:   3.1,
					Venue:   statistico.VenueEnum_AWAY,
				},
			},
			Visibility: statistico.VisibilityEnum_PRIVATE,
			StakingPlan: &statistico.StakingPlan{
				Name:  statistico.StakingPlanEnum_PERCENTAGE,
				Value: 2.5,
			},
		}

		st := mock.MatchedBy(func(s *strategy.Strategy) bool {
			res := []*strategy.ResultFilter{
				{
					Team:   "HOME_TEAM",
					Result: "WIN_DRAW",
					Games:  2,
					Venue:  "HOME_AWAY",
				},
			}

			stat := []*strategy.StatFilter{
				{
					Stat:    "GOALS",
					Team:    "HOME_TEAM",
					Action:  "AGAINST",
					Games:   4,
					Measure: "AVERAGE",
					Metric:  "GTE",
					Value:   3.1,
					Venue:   "AWAY",
				},
			}

			a := assert.New(t)

			a.Equal(r.GetName(), s.Name)
			a.Equal(r.GetDescription(), s.Description)
			a.Equal(r.GetMarket(), s.MarketName)
			a.Equal(r.GetRunner(), s.RunnerName)
			a.Equal(r.GetCompetitionIds(), s.CompetitionIDs)
			a.Equal(r.GetSide().String(), s.Side)
			a.Equal(r.GetVisibility().String(), s.Visibility)
			a.Equal("ACTIVE", s.Status)
			a.Equal(res, s.ResultFilters)
			a.Equal(stat, s.StatFilters)
			a.Equal(time.Unix(1616936636, 0), s.CreatedAt)
			a.Equal(time.Unix(1616936636, 0), s.UpdatedAt)
			return true
		})

		e := errors.New("error within writer")

		writer.On("Insert", st).Return(e)

		ctx := context.WithValue(context.Background(), "userID", "a5f04fd2-dfe7-41c1-af38-d490119705d8")

		_, err := service.SaveStrategy(ctx, r)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if e, ok := status.FromError(err); ok {
			assert.Equal(t, codes.Code(13), e.Code())
		}

		assert.Equal(t, "rpc error: code = Internal desc = internal server error", err.Error())

		writer.AssertExpectations(t)
	})
}

func TestStrategyService_ListUserStrategies(t *testing.T) {
	t.Run("streams user strategies fetched from strategy reader", func(t *testing.T) {
		t.Helper()

		writer := new(MockStrategyWriter)
		reader := new(MockStrategyReader)
		builder := new(MockStrategyBuilder)
		logger, _ := test.NewNullLogger()
		clock := clockwork.NewFakeClockAt(time.Unix(1616936636, 0))

		service := g.NewStrategyService(builder, writer, reader, logger, clock)

		stream := new(MockStrategyServer)

		r := statistico.ListUserStrategiesRequest{
			UserId:               "a5f04fd2-dfe7-41c1-af38-d490119705d8",
		}

		query := mock.MatchedBy(func(q *strategy.ReaderQuery) bool {
			assert.Equal(t, "a5f04fd2-dfe7-41c1-af38-d490119705d8", q.UserID.String())
			return true
		})

		strategies := []*strategy.Strategy{
			{
				ID:             uuid.New(),
				Name:           "Strategy One",
				Description:    "First Strategy",
				UserID:         uuid.MustParse("a5f04fd2-dfe7-41c1-af38-d490119705d8"),
				MarketName:     "MATCH_ODDS",
				RunnerName:     "Home",
				MinOdds:        nil,
				MaxOdds:        nil,
				CompetitionIDs: []uint64{8, 14},
				Side:           "BACK",
				Visibility:     "PUBLIC",
				Status:         "ACTIVE",
				StakingPlan:    strategy.StakingPlan{
					Name:   "PERCENTAGE",
					Number: 2.5,
				},
				ResultFilters:  []*strategy.ResultFilter{},
				StatFilters:    []*strategy.StatFilter{},
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		reader.On("Get", query).Return(strategies, nil)

		stream.On("Context").Return(context.WithValue(context.Background(), "userID", "a5f04fd2-dfe7-41c1-af38-d490119705d8"))
		stream.On("Send", mock.AnythingOfType("*statistico.Strategy")).Once().Return(nil)

		err := service.ListUserStrategies(&r, stream)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		stream.AssertExpectations(t)
		reader.AssertExpectations(t)
	})
}

type MockStrategyBuilder struct {
	mock.Mock
}

func (m *MockStrategyBuilder) Build(ctx context.Context, q *strategy.BuilderQuery) <-chan *strategy.Trade {
	args := m.Called(ctx, q)
	return args.Get(0).(<-chan *strategy.Trade)
}

type MockStrategyWriter struct {
	mock.Mock
}

func (m *MockStrategyWriter) Insert(s *strategy.Strategy) error {
	args := m.Called(s)
	return args.Error(0)
}

type MockStrategyReader struct {
	mock.Mock
}

func (m *MockStrategyReader) Get(q *strategy.ReaderQuery) ([]*strategy.Strategy, error) {
	args := m.Called(q)
	return args.Get(0).([]*strategy.Strategy), args.Error(1)
}

type MockStrategyBuildServer struct {
	mock.Mock
	grpc.ServerStream
}

func (m *MockStrategyBuildServer) Send(s *statistico.StrategyTrade) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockStrategyBuildServer) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

type MockStrategyServer struct {
	mock.Mock
	grpc.ServerStream
}

func (m *MockStrategyServer) Send(s *statistico.Strategy) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockStrategyServer) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}


func tradeChannel(trades []*strategy.Trade) <-chan *strategy.Trade {
	ch := make(chan *strategy.Trade, len(trades))

	for _, m := range trades {
		ch <- m
	}

	close(ch)

	return ch
}
