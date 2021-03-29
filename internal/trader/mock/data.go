package mock

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/mock"
)

type FixtureClient struct {
	mock.Mock
}

func (m *FixtureClient) ByID(ctx context.Context, fixtureID uint64) (*statistico.Fixture, error) {
	args := m.Called(ctx, fixtureID)
	return args.Get(0).(*statistico.Fixture), args.Error(1)
}

func (m *FixtureClient) Search(ctx context.Context, req *statistico.FixtureSearchRequest) ([]*statistico.Fixture, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*statistico.Fixture), args.Error(1)
}

type ResultClient struct {
	mock.Mock
}

func (m *ResultClient) ByID(ctx context.Context, fixtureID uint64) (*statistico.Result, error) {
	args := m.Called(ctx, fixtureID)
	return args.Get(0).(*statistico.Result), args.Error(1)
}

func (m *ResultClient) ByTeam(ctx context.Context, req *statistico.TeamResultRequest) ([]*statistico.Result, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*statistico.Result), args.Error(1)
}
