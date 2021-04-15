package trade

import (
	"context"
	"github.com/statistico/statistico-trader/internal/trader/auth"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
)

type manager struct {
	factory exchange.ClientFactory
	users   auth.UserService
	placer  Placer
}

func (m *manager) Manage(ctx context.Context, t *Ticket, s *strategy.Strategy) error {
	user, err := m.users.ByID(s.UserID)

	if err != nil {
		return err
	}

	client, err := m.factory.Create(t.Exchange, user.BetFairUserName, user.BetFairPassword, user.BetFairKey)

	if err != nil {
		return err
	}

	_, err = m.placer.PlaceTrade(ctx, client, t, s)

	if err != nil {
		return err
	}

	return nil
}

func NewManager(u auth.UserService, p Placer) Manager {
	return &manager{
		users:  u,
		placer: p,
	}
}
