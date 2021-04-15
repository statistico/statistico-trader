package trade_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/auth"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/statistico/statistico-trader/internal/trader/trade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestManager_Manage(t *testing.T) {
	t.Run("fetches user and places trade via trade.Placer", func(t *testing.T) {
		t.Helper()

		factory := new(MockExchangeClientFactory)
		users := new(MockUserService)
		placer := new(MockTradePlacer)
		manager := trade.NewManager(factory, users, placer)

		ctx := context.Background()
		s := strategy.Strategy{UserID: uuid.MustParse("794fe24b-6a8f-4fe7-b235-05cef412b80e")}
		ticket := trade.Ticket{Exchange: "betfair"}

		user := auth.User{
			ID:              s.ID,
			Email:           "joe@email.com",
			BetFairUserName: "joe",
			BetFairPassword: "password",
			BetFairKey:      "key-123",
		}

		users.On("ByID", s.UserID).Return(&user, nil)

		client := new(MockExchangeClient)

		factory.On("Create", "betfair", "joe", "password", "key-123").Return(client, nil)

		placer.On("PlaceTrade", ctx, client, &ticket, &s).Return(&trade.Trade{}, nil)

		err := manager.Manage(ctx, &ticket, &s)

		if err != nil {
			t.Fatalf("Expected nil, got %+v", err)
		}

		users.AssertExpectations(t)
		factory.AssertExpectations(t)
		placer.AssertExpectations(t)
	})

	t.Run("returns error if error returned by auth.UserService", func(t *testing.T) {
		t.Helper()

		factory := new(MockExchangeClientFactory)
		users := new(MockUserService)
		placer := new(MockTradePlacer)
		manager := trade.NewManager(factory, users, placer)

		ctx := context.Background()
		s := strategy.Strategy{UserID: uuid.MustParse("794fe24b-6a8f-4fe7-b235-05cef412b80e")}
		ticket := trade.Ticket{Exchange: "betfair"}

		users.On("ByID", s.UserID).Return(&auth.User{}, errors.New("user service error"))

		factory.AssertNotCalled(t, "Create")
		placer.AssertNotCalled(t, "PlaceTrade")

		err := manager.Manage(ctx, &ticket, &s)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "user service error", err.Error())
	})

	t.Run("returns error if error returned by exchange.ClientFactory", func(t *testing.T) {
		t.Helper()

		factory := new(MockExchangeClientFactory)
		users := new(MockUserService)
		placer := new(MockTradePlacer)
		manager := trade.NewManager(factory, users, placer)

		ctx := context.Background()
		s := strategy.Strategy{UserID: uuid.MustParse("794fe24b-6a8f-4fe7-b235-05cef412b80e")}
		ticket := trade.Ticket{Exchange: "betfair"}

		user := auth.User{
			ID:              s.ID,
			Email:           "joe@email.com",
			BetFairUserName: "joe",
			BetFairPassword: "password",
			BetFairKey:      "key-123",
		}

		users.On("ByID", s.UserID).Return(&user, nil)

		client := new(MockExchangeClient)

		factory.On("Create", "betfair", "joe", "password", "key-123").Return(client, errors.New("factory error"))

		placer.AssertNotCalled(t, "PlaceTrade")

		err := manager.Manage(ctx, &ticket, &s)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "factory error", err.Error())
	})

	t.Run("returns nil if trade.DuplicationError returned by trade.Placer", func(t *testing.T) {
		t.Helper()

		factory := new(MockExchangeClientFactory)
		users := new(MockUserService)
		placer := new(MockTradePlacer)
		manager := trade.NewManager(factory, users, placer)

		ctx := context.Background()
		s := strategy.Strategy{UserID: uuid.MustParse("794fe24b-6a8f-4fe7-b235-05cef412b80e")}
		ticket := trade.Ticket{Exchange: "betfair"}

		user := auth.User{
			ID:              s.ID,
			Email:           "joe@email.com",
			BetFairUserName: "joe",
			BetFairPassword: "password",
			BetFairKey:      "key-123",
		}

		users.On("ByID", s.UserID).Return(&user, nil)

		client := new(MockExchangeClient)

		factory.On("Create", "betfair", "joe", "password", "key-123").Return(client, nil)

		placer.On("PlaceTrade", ctx, client, &ticket, &s).Return(&trade.Trade{}, &trade.DuplicationError{})

		err := manager.Manage(ctx, &ticket, &s)

		if err != nil {
			t.Fatalf("Expected nil, got %+v", err)
		}

		users.AssertExpectations(t)
		factory.AssertExpectations(t)
		placer.AssertExpectations(t)
	})

	t.Run("returns error if non trade.DuplicationError returned by trade.Placer", func(t *testing.T) {
		t.Helper()

		factory := new(MockExchangeClientFactory)
		users := new(MockUserService)
		placer := new(MockTradePlacer)
		manager := trade.NewManager(factory, users, placer)

		ctx := context.Background()
		s := strategy.Strategy{UserID: uuid.MustParse("794fe24b-6a8f-4fe7-b235-05cef412b80e")}
		ticket := trade.Ticket{Exchange: "betfair"}

		user := auth.User{
			ID:              s.ID,
			Email:           "joe@email.com",
			BetFairUserName: "joe",
			BetFairPassword: "password",
			BetFairKey:      "key-123",
		}

		users.On("ByID", s.UserID).Return(&user, nil)

		client := new(MockExchangeClient)

		factory.On("Create", "betfair", "joe", "password", "key-123").Return(client, nil)

		placer.On("PlaceTrade", ctx, client, &ticket, &s).Return(&trade.Trade{}, errors.New("placer error"))

		err := manager.Manage(ctx, &ticket, &s)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "placer error", err.Error())

		users.AssertExpectations(t)
		factory.AssertExpectations(t)
		placer.AssertExpectations(t)
	})
}

type MockExchangeClientFactory struct {
	mock.Mock
}

func (m *MockExchangeClientFactory) Create(e, user, password, key string) (exchange.Client, error) {
	args := m.Called(e, user, password, key)
	return args.Get(0).(exchange.Client), args.Error(1)
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
