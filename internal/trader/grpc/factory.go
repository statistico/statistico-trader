package grpc

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func strategyFromRequest(r *statistico.SaveStrategyRequest, t time.Time) (*trader.Strategy, error) {
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

	userID, err := uuid.Parse(r.GetUserId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error parsing user ID: %s", err.Error())
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

func parseStakingPlan(s *statistico.StakingPlan) (trader.StakingPlan, error) {
	if s.Name.String() != "PERCENTAGE" {
		return trader.StakingPlan{}, fmt.Errorf("staking plan '%s' is not supported", s.Name)
	}

	if s.Value <= 0 {
		return trader.StakingPlan{}, errors.New("staking plan must be greater than zero")
	}

	return trader.StakingPlan{
		Name:  s.Name.String(),
		Number: s.Value,
	}, nil
}
