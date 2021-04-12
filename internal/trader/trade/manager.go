package trade

import (
	"context"
	"github.com/statistico/statistico-trader/internal/trader/auth"
	"github.com/statistico/statistico-trader/internal/trader/market"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
)

type manager struct {
	users   auth.UserService
	placer  Placer
}

func (m *manager) Manage(ctx context.Context, r *market.Runner, s *strategy.Strategy) error {
	
}