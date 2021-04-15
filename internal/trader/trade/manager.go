package trade

import (
	"context"
	betfair "github.com/statistico/statistico-betfair-go-client"
	"github.com/statistico/statistico-trader/internal/trader/auth"
	betfair2 "github.com/statistico/statistico-trader/internal/trader/exchange/betfair"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"net/http"
)

type manager struct {
	users   auth.UserService
	placer  Placer
}

func (m *manager) Manage(ctx context.Context, t *Ticket, s *strategy.Strategy) error {
	user, err := m.users.ByID(s.UserID)

	if err != nil {
		return err
	}

	// Move logic into ExchangeClientFactory
	c := betfair.NewClient(&http.Client{}, betfair.InteractiveCredentials{
		Username: user.BetFairUserName,
		Password: user.BetFairPassword,
		Key:      user.BetFairKey,
	})

	client := betfair2.NewExchangeClient(c)

	// Will send notification to user with returned trade
	_, err = m.placer.PlaceTrade(ctx, client, t, s)

	switch e := err.(type) {
	case *DuplicationError:
		return nil
	case nil:
		return nil
	default:
		return e
	}
}

func NewManager(u auth.UserService, p Placer) Manager {
	return &manager{
		users:  u,
		placer: p,
	}
}
