package trade

import (
	"context"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/market"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"math"
)

type manager struct {
	reader Reader
	writer Writer
	clock clockwork.Clock
}

func (m *manager) PlaceTrade(ctx context.Context, c exchange.Client, r *market.Runner, s strategy.Strategy) (*Trade, error) {
	exists, err := m.reader.Exists(r.MarketName, r.RunnerName, r.EventID, s.ID)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, &DuplicationError{
			market:     r.MarketName,
			runner:     r.RunnerName,
			eventID:    r.EventID,
			strategyID: s.ID,
		}
	}

	account, err := c.Account(ctx)

	if err != nil {
		return nil, &ExchangeError{err: err}
	}

	stake, err := calculateStake(account, s.StakingPlan)

	if err != nil {
		return nil, err
	}

	if stake <= 0 {
		return nil, &InvalidBalanceError{
			market:     r.MarketName,
			runner:     r.RunnerName,
			eventID:    r.EventID,
			strategyID: s.ID,
			balance:    stake,
		}
	}

	ticket := exchange.TradeTicket{
		MarketID: r.MarketID,
		RunnerID: r.RunnerID,
		Price:    r.Price.Value,
		Stake:    stake,
		Side:     r.Price.Side,
	}

	response, err := c.PlaceTrade(ctx, &ticket)

	if err != nil {
		return nil, &ExchangeError{err: err}
	}

	t := Trade{
		ID:          uuid.New(),
		StrategyID:  s.ID,
		Exchange:    response.Exchange,
		ExchangeRef: response.Reference,
		Market:      r.MarketName,
		Runner:      r.RunnerName,
		Price:       ticket.Price,
		Stake:       ticket.Stake,
		EventID:     r.EventID,
		EventDate:   r.EventDate,
		Side:        ticket.Side,
		Result:      InPlay,
		Timestamp:   m.clock.Now(),
	}

	if err := m.writer.Insert(&t); err != nil {
		return &t, err
	}

	return &t, nil
}

func calculateStake(account *exchange.Account, plan strategy.StakingPlan) (float32, error) {
	total := float64(account.Balance) + math.Abs(float64(account.Exposure))

	if total <= 0 {
		return float32(total), nil
	}

	if plan.Name != strategy.PercentageStakingPlan {
		return 0, nil
	}

	return float32(total) * plan.Number, nil
}

func NewManager(r Reader, w Writer, c clockwork.Clock) Manager {
	return &manager{
		reader: r,
		writer: w,
		clock:  c,
	}
}
