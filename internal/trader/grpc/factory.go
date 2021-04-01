package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func strategyFromRequest(ctx context.Context, r *statistico.SaveStrategyRequest, t time.Time) (*trader.Strategy, error) {
	strategy := trader.Strategy{
		ID:             uuid.New(),
		Name:           r.GetName(),
		Description:    r.GetDescription(),
		MarketName:     r.GetMarket(),
		RunnerName:     r.GetRunner(),
		CompetitionIDs: r.GetCompetitionIds(),
		Side:           r.GetSide().String(),
		Visibility:     r.GetVisibility().String(),
		Status:         "ACTIVE",
		ResultFilters:  transformResultFilters(r.ResultFilters),
		StatFilters:    transformStatFilters(r.StatFilters),
		CreatedAt:      t,
		UpdatedAt:      t,
	}

	userID, err := uuid.Parse(fmt.Sprintf("%v", ctx.Value("userID")))

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user id provided is not a valid uuid string: %s", err.Error())
	}

	strategy.UserID = userID

	plan, err := parseStakingPlan(r.StakingPlan)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	strategy.StakingPlan = plan

	if r.GetMinOdds() == nil && r.GetMaxOdds() == nil {
		return nil, status.Error(codes.InvalidArgument, "Min and max odds cannot both be nil")
	}

	if r.GetMinOdds() != nil {
		strategy.MinOdds = &r.GetMinOdds().Value
	}

	if r.GetMaxOdds() != nil {
		strategy.MaxOdds = &r.GetMaxOdds().Value
	}

	return &strategy, nil
}

func strategyReaderQuery(ctx context.Context, r *statistico.ListUserStrategiesRequest) (*trader.StrategyReaderQuery, error) {
	userID, err := uuid.Parse(r.GetUserId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user id provided is not a valid uuid string: %s", err.Error())
	}

	query := trader.StrategyReaderQuery{
		UserID:     &userID,
	}

	if userID.String() != ctx.Value("userID") {
		visibility := "PUBLIC"
		query.Visibility = &visibility
	}

	return &query, nil
}

func convertToStatisticoStrategy(s *trader.Strategy) *statistico.Strategy {
	st := statistico.Strategy{
		Id:             s.ID.String(),
		Name:           s.Name,
		Description:    s.Description,
		UserId:         s.UserID.String(),
		Market:         s.MarketName,
		Runner:         s.RunnerName,
		CompetitionIds: s.CompetitionIDs,
		Side:           statistico.SideEnum(statistico.SideEnum_value[s.Side]),
		Visibility:     statistico.VisibilityEnum(statistico.VisibilityEnum_value[s.Visibility]),
		Status:         statistico.StrategyStatusEnum(statistico.StrategyStatusEnum_value[s.Status]),
		ResultFilters:  convertResultFilters(s.ResultFilters),
		StatFilters:    convertStatFilters(s.StatFilters),
		CreatedAt:      timestamppb.New(s.CreatedAt),
		UpdatedAt:      timestamppb.New(s.UpdatedAt),
	}

	if s.MinOdds != nil {
		st.MinOdds = &wrappers.FloatValue{Value: *s.MinOdds}
	}

	if s.MaxOdds != nil {
		st.MaxOdds = &wrappers.FloatValue{Value: *s.MaxOdds}
	}

	return &st
}

func convertResultFilters(f []*trader.ResultFilter) []*statistico.ResultFilter {
	filters := []*statistico.ResultFilter{}

	for _, ft := range f {
		srf := &statistico.ResultFilter{
			Team:                 statistico.TeamEnum(statistico.TeamEnum_value[ft.Team]),
			Result:               statistico.ResultEnum(statistico.ResultEnum_value[ft.Result]),
			Games:                uint32(ft.Games),
			Venue:                statistico.VenueEnum(statistico.VenueEnum_value[ft.Venue]),
		}

		filters = append(filters, srf)
	}

	return filters
}

func convertStatFilters(f []*trader.StatFilter) []*statistico.StatFilter {
	filters := []*statistico.StatFilter{}

	for _, ft := range f {
		sst := &statistico.StatFilter{
			Stat:                 statistico.StatEnum(statistico.StatEnum_value[ft.Stat]),
			Team:                 statistico.TeamEnum(statistico.TeamEnum_value[ft.Team]),
			Action:               statistico.ActionEnum(statistico.ActionEnum_value[ft.Action]),
			Games:                uint32(ft.Games),
			Measure:              statistico.MeasureEnum(statistico.MeasureEnum_value[ft.Measure]),
			Metric:               statistico.MetricEnum(statistico.MetricEnum_value[ft.Metric]),
			Value:                ft.Value,
			Venue:                statistico.VenueEnum(statistico.VenueEnum_value[ft.Venue]),
		}

		filters = append(filters, sst)
	}

	return filters
}

func parseStakingPlan(s *statistico.StakingPlan) (trader.StakingPlan, error) {
	if s.Name.String() != "PERCENTAGE" {
		return trader.StakingPlan{}, fmt.Errorf("staking plan '%s' is not supported", s.Name)
	}

	if s.Value <= 0 {
		return trader.StakingPlan{}, errors.New("staking plan must be greater than zero")
	}

	return trader.StakingPlan{
		Name:   s.Name.String(),
		Number: s.Value,
	}, nil
}
