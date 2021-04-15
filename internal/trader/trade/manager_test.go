package trade_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/auth"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/statistico/statistico-trader/internal/trader/trade"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestManager_Manage(t *testing.T) {
	t.Run("fetches user and places trade via trade.Placer", func(t *testing.T) {
		t.Helper()
		
		users := new(MockUserService)
		placer := new(MockTradePlacer)
		
		manager := trade.NewManager(users, placer)

		ctx := context.Background()
		s := strategy.Strategy{ID: uuid.MustParse("794fe24b-6a8f-4fe7-b235-05cef412b80e")}
		ticket := trade.Ticket{}

		user := auth.User{
			ID:              s.ID,
			Email:           "joe@email.com",
			BetFairUserName: "joe",
			BetFairPassword: "password",
			BetFairKey:      "key-123",
		}
	})
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) ByID(userID uuid.UUID) (*auth.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*auth.User), args.Error(1)
}

type MockTradePlacer struct {
	mock.Mock
}

func (m *MockTradePlacer) PlaceTrade(ctx context.Context, c exchange.Client, t *trade.Ticket, s *strategy.Strategy) (*trade.Trade, error) {
	args := m.Called(ctx, c, t, s)
	return args.Get(0).(*trade.Trade), args.Error(1)
}
