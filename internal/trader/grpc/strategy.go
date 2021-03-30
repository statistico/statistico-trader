package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-odds-warehouse-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StrategyService struct {
	writer     trader.StrategyWriter
	oddsClient statisticooddswarehouse.MarketClient
	finder     TradeFinder
	logger     *logrus.Logger
	clock      clockwork.Clock
	statistico.UnimplementedStrategyServiceServer
}

func (s *StrategyService) BuildStrategy(r *statistico.BuildStrategyRequest, stream statistico.StrategyService_BuildStrategyServer) error {
	req := statistico.MarketRunnerRequest{
		Market:         r.GetMarket(),
		Runner:         r.GetRunner(),
		Line:           r.GetLine(),
		Side:           r.GetSide(),
		MinOdds:        r.GetMinOdds(),
		MaxOdds:        r.GetMaxOdds(),
		CompetitionIds: r.GetCompetitionIds(),
		SeasonIds:      r.GetSeasonIds(),
		DateFrom:       r.GetDateFrom(),
		DateTo:         r.GetDateTo(),
	}

	ctx := context.Background()

	markets, errCh := s.oddsClient.MarketRunnerSearch(ctx, &req)

	query := TradeQuery{
		Markets:       markets,
		RunnerFilters: transformResultFilters(r.ResultFilters),
		StatFilters:   transformStatFilters(r.StatFilters),
	}

	trades := s.finder.Find(ctx, &query)

	for t := range trades {
		if err := stream.Send(t); err != nil {
			s.logger.Errorf("error streaming market runner back to client: %s", err.Error())
		}
	}

	err := <-errCh

	if err != nil {
		s.logger.Errorf("error fetching market runners from odds warehouse: %s", err.Error())
	}

	return err
}

func (s *StrategyService) SaveStrategy(ctx context.Context, r *statistico.SaveStrategyRequest) (*statistico.Strategy, error) {
	strategy, err := strategyFromRequest(r, s.clock.Now())

	if err != nil {
		return nil, err
	}

	err = s.writer.Insert(strategy)

	if err != nil {
		if de, ok := err.(*errors.DuplicationError); ok {
			return nil, status.Error(codes.AlreadyExists, de.Error())
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &statistico.Strategy{
		Id:             strategy.ID.String(),
		Name:           strategy.Name,
		Description:    strategy.Description,
		UserId:         strategy.UserID.String(),
		Market:         strategy.MarketName,
		Runner:         strategy.RunnerName,
		MinOdds:        &wrappers.FloatValue{Value: *strategy.MinOdds},
		MaxOdds:        &wrappers.FloatValue{Value: *strategy.MaxOdds},
		CompetitionIds: strategy.CompetitionIDs,
		Side:           statistico.SideEnum(statistico.SideEnum_value[strategy.Side]),
		Visibility:     statistico.VisibilityEnum(statistico.VisibilityEnum_value[strategy.Visibility]),
		Status:         statistico.StrategyStatusEnum_ACTIVE,
		ResultFilters:  r.ResultFilters,
		StatFilters:    r.StatFilters,
		CreatedAt:      timestamppb.New(strategy.CreatedAt),
		UpdatedAt:      timestamppb.New(strategy.UpdatedAt),
	}, nil
}

func NewStrategyService(
	w trader.StrategyWriter,
	c statisticooddswarehouse.MarketClient,
	f TradeFinder,
	l *logrus.Logger,
	cl clockwork.Clock,
) *StrategyService {
	return &StrategyService{
		writer:     w,
		oddsClient: c,
		finder:     f,
		logger:     l,
		clock:      cl,
	}
}
