package market

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-trader/internal/trader/queue"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/statistico/statistico-trader/internal/trader/trade"
	"sync"
)

type Handler interface {
	HandleEventMarket(ctx context.Context, e *queue.EventMarket)
}

type handler struct {
	finder  strategy.Finder
	manager trade.Manager
	logger  *logrus.Logger
}

func (h *handler) HandleEventMarket(ctx context.Context, e *queue.EventMarket) {
	wg := sync.WaitGroup{}

	for _, runner := range e.Runners {
		wg.Add(2)
		go h.handleBackRunner(ctx, e, runner, &wg)
		go h.handleLayRunner(ctx, e, runner, &wg)
	}

	wg.Wait()
}

func (h *handler) handleBackRunner(ctx context.Context, e *queue.EventMarket, r *queue.Runner, wg *sync.WaitGroup) {
	if len(r.BackPrices) < 1 {
		wg.Done()
		return
	}

	price := r.BackPrices[0]

	t := trade.Ticket{
		MarketID:      e.ID,
		MarketName:    e.Name,
		RunnerID:      r.ID,
		RunnerName:    r.Name,
		EventID:       e.EventID,
		CompetitionID: e.CompetitionID,
		SeasonID:      e.SeasonID,
		EventDate:     e.EventDate,
		Exchange:      e.Exchange,
		Price:         trade.TicketPrice{
			Value:     price.Price,
			Size:      price.Size,
			Side:      strategy.Back,
		},
	}

	h.handleRunner(ctx, &t, wg)

	wg.Done()
}

func (h *handler) handleLayRunner(ctx context.Context, e *queue.EventMarket, r *queue.Runner, wg *sync.WaitGroup) {
	if len(r.LayPrices) < 1 {
		wg.Done()
		return
	}

	price := r.LayPrices[0]

	t := trade.Ticket{
		MarketID:      e.ID,
		MarketName:    e.Name,
		RunnerID:      r.ID,
		RunnerName:    r.Name,
		EventID:       e.EventID,
		CompetitionID: e.CompetitionID,
		SeasonID:      e.SeasonID,
		EventDate:     e.EventDate,
		Exchange:      e.Exchange,
		Price:         trade.TicketPrice{
			Value:     price.Price,
			Size:      price.Size,
			Side:      strategy.Back,
		},
	}

	h.handleRunner(ctx, &t, wg)

	wg.Done()
}

func (h *handler) handleRunner(ctx context.Context, t *trade.Ticket, wg *sync.WaitGroup) {
	query := strategy.FinderQuery{
		MarketName:    t.MarketName,
		RunnerName:    t.RunnerName,
		EventID:       t.EventID,
		CompetitionID: t.CompetitionID,
		Price:         t.Price.Value,
		Side:          t.Price.Side,
		Status:        strategy.Active,
	}

	st := h.finder.FindMatchingStrategies(ctx, &query)

	for s := range st {
		wg.Add(1)

		go func(s *strategy.Strategy) {
			if err := h.manager.Manage(ctx, t, s); err != nil {
				h.logger.Errorf("error managing trade for strategy %s and market %s: %+v", s.ID, t.MarketName, err)
			}
			wg.Done()
		}(s)
	}
}

func NewHandler(f strategy.Finder, t trade.Manager, l *logrus.Logger) Handler {
	return &handler{
		finder:  f,
		manager: t,
		logger:  l,
	}
}
