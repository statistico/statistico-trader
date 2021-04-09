package market

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-trader/internal/trader"
	"github.com/statistico/statistico-trader/internal/trader/classify"
	"github.com/statistico/statistico-trader/internal/trader/queue"
	"sync"
)

const (
	Active = "ACTIVE"
	Back = "BACK"
	Lay = "LAY"
)

type Handler interface {
	HandleEventMarket(e *queue.EventMarket) error
}

type handler struct {
	reader  trader.StrategyReader
	// Add strategy service / trade service here
	logger  logrus.Logger
}

func (h *handler) HandleEventMarket(e *queue.EventMarket) error {
	wg := sync.WaitGroup{}

	for _, runner := range e.Runners {
		a := Active
		b := Back

		query := trader.StrategyReaderQuery{
			Market:        &e.Name,
			Runner:        &runner.Name,
			Price:         &runner.BackPrices[0].Size,
			CompetitionID: &e.CompetitionID,
			Side:          &b,
			Status:        &a,
		}

		st, err := h.reader.Get(&query)

		if err != nil {
			h.logger.Error("error fetching strategies in market handler: %+v", err)
		}

		// Loop over strategies and pass individual strategy to Trade Service / Strategy Handler in go routine

		// Repeat for LAY market/strategy
	}
}
