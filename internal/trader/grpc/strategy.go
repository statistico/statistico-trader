package grpc

import (
	"context"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/errors"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StrategyService struct {
	builder    strategy.Builder
	reader     strategy.Reader
	writer     strategy.Writer
	logger     *logrus.Logger
	clock      clockwork.Clock
	statistico.UnimplementedStrategyServiceServer
}

func (s *StrategyService) HealthCheck(ctx context.Context, r *statistico.HealthCheckRequest) (*statistico.HealthCheckResponse, error) {
	s.logger.Errorf("Request received")
	return &statistico.HealthCheckResponse{Message: "HealthCheck OK from statistico-trader"}, nil
}

func (s *StrategyService) BuildStrategy(r *statistico.BuildStrategyRequest, stream statistico.StrategyService_BuildStrategyServer) error {
	query := strategy.BuilderQuery{
		Market:         r.GetMarket(),
		Runner:         r.GetRunner(),
		Line:           r.GetLine(),
		Side:           r.GetSide().String(),
		CompetitionIDs: r.GetCompetitionIds(),
		SeasonIDs:      r.GetSeasonIds(),
		ResultFilters:  transformResultFilters(r.GetResultFilters()),
		StatFilters:    transformStatFilters(r.GetStatFilters()),
	}

	if r.GetMinOdds() != nil {
		query.MinOdds = &r.GetMinOdds().Value
	}

	if r.GetMaxOdds() != nil {
		query.MaxOdds = &r.GetMaxOdds().Value
	}

	ch := s.builder.Build(stream.Context(), &query)

	for t := range ch {
		if err := stream.Send(transformStrategyTrade(t)); err != nil {
			s.logger.Errorf("error streaming strategy trade back to client: %s", err.Error())
		}
	}

	return nil
}

func (s *StrategyService) SaveStrategy(ctx context.Context, r *statistico.SaveStrategyRequest) (*statistico.Strategy, error) {
	st, err := strategyFromRequest(ctx, r, s.clock.Now())

	if err != nil {
		s.logger.Errorf("error parsing strategy request %+v", err)
		return nil, err
	}

	err = s.writer.Insert(st)

	if err != nil {
		s.logger.Errorf("error saving strategy %+v", err)

		if de, ok := err.(*errors.DuplicationError); ok {
			return nil, status.Error(codes.AlreadyExists, de.Error())
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return convertToStatisticoStrategy(st), nil
}

func (s *StrategyService) ListUserStrategies(r *statistico.ListUserStrategiesRequest, stream statistico.StrategyService_ListUserStrategiesServer) error {
	s.logger.Errorf("Request received")
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

	s.logger.Errorf("Response returned")

	return nil
}

func NewStrategyService(
	b strategy.Builder,
	w strategy.Writer,
	r strategy.Reader,
	l *logrus.Logger,
	cl clockwork.Clock,
) *StrategyService {
	return &StrategyService{
		builder:    b,
		writer:     w,
		reader:     r,
		logger:     l,
		clock:      cl,
	}
}
