package market

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-trader/internal/trader"
	"github.com/statistico/statistico-trader/internal/trader/queue"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"sync"
	"time"
)

type Handler interface {
	HandleEventMarket(e *queue.EventMarket) error
}

type handler struct {
	finder  strategy.Finder
	// Add strategy service / trade service here
	logger  logrus.Logger
}

func (h *handler) HandleEventMarket(e *queue.EventMarket) error {
	wg := sync.WaitGroup{}
	ch := make(chan *strategy.Strategy, 100)

	for _, runner := range e.Runners {
		wg.Add(1)

		go h.handleBackRunner(e, runner, ch, &wg)
		
		// Loop over strategies and pass individual strategy to Trade Service / Strategy Handler in go routine

		// Repeat for LAY market/strategy
	}

	wg.Wait()
}

func (h *handler) handleBackRunner(e *queue.EventMarket, r *queue.Runner, ch chan<- *strategy.Strategy, wg *sync.WaitGroup) {
	if len(r.BackPrices) < 1 {
		wg.Done()
		return
	}

	price := r.BackPrices[0]

	mk := Runner{
		MarketID:      e.ID,
		MarketName:    e.Name,
		RunnerID:      r.ID,
		RunnerName:    r.Name,
		EventID:       e.EventID,
		CompetitionID: e.CompetitionID,
		SeasonID:      e.SeasonID,
		EventDate:     e.EventDate,
		Exchange:      e.Exchange,
		Price:         Price{
			Value:     price.Price,
			Size:      price.Size,
			Side:      strategy.Back,
		},
	}

	h.finder.FindMatchingStrategies(context.Background(), &mk, ch)

	wg.Done()
}