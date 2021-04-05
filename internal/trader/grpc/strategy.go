package grpc

import (
	"context"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-odds-warehouse-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StrategyService struct {
	writer     trader.StrategyWriter
	reader     trader.StrategyReader
	oddsClient statisticooddswarehouse.MarketClient
	logger     *logrus.Logger
	clock      clockwork.Clock
	statistico.UnimplementedStrategyServiceServer
}

func (s *StrategyService) BuildStrategy(r *statistico.BuildStrategyRequest, stream statistico.StrategyService_BuildStrategyServer) error {
	//req := statistico.MarketRunnerRequest{
	//	Market:         r.GetMarket(),
	//	Runner:         r.GetRunner(),
	//	Line:           r.GetLine(),
	//	Side:           r.GetSide(),
	//	MinOdds:        r.GetMinOdds(),
	//	MaxOdds:        r.GetMaxOdds(),
	//	CompetitionIds: r.GetCompetitionIds(),
	//	SeasonIds:      r.GetSeasonIds(),
	//	DateFrom:       r.GetDateFrom(),
	//	DateTo:         r.GetDateTo(),
	//}
	//
	//ctx := context.Background()
	//
	//markets, errCh := s.oddsClient.MarketRunnerSearch(ctx, &req)
	//
	//query := TradeQuery{
	//	Markets:       markets,
	//	ResultFilters: transformResultFilters(r.ResultFilters),
	//	StatFilters:   transformStatFilters(r.StatFilters),
	//}
	//
	//trades := s.finder.Find(ctx, &query)
	//
	//for t := range trades {
	//	if err := stream.Send(t); err != nil {
	//		s.logger.Errorf("error streaming market runner back to client: %s", err.Error())
	//	}
	//}
	//
	//err := <-errCh
	//
	//if err != nil {
	//	s.logger.Errorf("error fetching market runners from odds warehouse: %s", err.Error())
	//}

	return nil
}

func (s *StrategyService) SaveStrategy(ctx context.Context, r *statistico.SaveStrategyRequest) (*statistico.Strategy, error) {
	strategy, err := strategyFromRequest(ctx, r, s.clock.Now())

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

	return convertToStatisticoStrategy(strategy), nil
}

func (s *StrategyService) ListUserStrategies(r *statistico.ListUserStrategiesRequest, stream statistico.StrategyService_ListUserStrategiesServer) error {
	query, err :=  strategyReaderQuery(stream.Context(), r)

	if err != nil {
		return err
	}

	strategies, err := s.reader.Get(query)

	if err != nil {
		s.logger.Errorf("error fetching strategies from reader: %s", err.Error())
		return status.Error(codes.Internal, "internal server error")
	}

	for _, st := range strategies {
		if err := stream.Send(convertToStatisticoStrategy(st)); err != nil {
			s.logger.Errorf("error streaming strategy back to client: %s", err.Error())
		}
	}

	return nil
}

func NewStrategyService(
	w trader.StrategyWriter,
	r trader.StrategyReader,
	c statisticooddswarehouse.MarketClient,
	l *logrus.Logger,
	cl clockwork.Clock,
) *StrategyService {
	return &StrategyService{
		writer:     w,
		reader:     r,
		oddsClient: c,
		logger:     l,
		clock:      cl,
	}
}
