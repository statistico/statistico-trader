package trade

import (
	"context"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-trader/internal/trader/exchange"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"math"
)

type placer struct {
	reader Reader
	writer Writer
	clock clockwork.Clock
}

func (p *placer) PlaceTrade(ctx context.Context, c exchange.Client, t *Ticket, s *strategy.Strategy) (*Trade, error) {
	exists, err := p.reader.Exists(t.MarketName, t.RunnerName, t.EventID, s.ID)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, &DuplicationError{
			market:     t.MarketName,
			runner:     t.RunnerName,
			eventID:    t.EventID,
			strategyID: s.ID,
		}
	}

	account, err := c.Account(ctx)

	if err != nil {
		return nil, &ExchangeError{err: err}
	}

	stake := calculateStake(account, s.StakingPlan)

	if stake <= 0 {
		return nil, &InvalidBalanceError{
			market:     t.MarketName,
			runner:     t.RunnerName,
			eventID:    t.EventID,
			strategyID: s.ID,
			balance:    stake,
		}
	}

	ticket := exchange.TradeTicket{
		MarketID: t.MarketID,
		RunnerID: t.RunnerID,
		Price:    t.Price.Value,
		Stake:    stake,
		Side:     t.Price.Side,
	}

	response, err := c.PlaceTrade(ctx, &ticket)

	if err != nil {
		return nil, &ExchangeError{err: err}
	}

	tr := Trade{
		ID:          uuid.New(),
		StrategyID:  s.ID,
		Exchange:    response.Exchange,
		ExchangeRef: response.Reference,
		Market:      t.MarketName,
		Runner:      t.RunnerName,
		Price:       ticket.Price,
		Stake:       ticket.Stake,
		EventID:     t.EventID,
		EventDate:   t.EventDate,
		Side:        ticket.Side,
		Result:      InPlay,
		Timestamp:   p.clock.Now(),
	}

	if err := p.writer.Insert(&tr); err != nil {
		return &tr, err
	}

	return &tr, nil
}

func calculateStake(account *exchange.Account, plan strategy.StakingPlan) float32 {
	if account.Balance == 0 {
		return 0
	}

	total := float64(account.Balance) + math.Abs(float64(account.Exposure))
	
	if plan.Name != strategy.PercentageStakingPlan {
		return 0
	}

	return (float32(total) / 100) * plan.Number
}

func NewPlacer(r Reader, w Writer, c clockwork.Clock) Placer {
	return &placer{
		reader: r,
		writer: w,
		clock:  c,
	}
}
