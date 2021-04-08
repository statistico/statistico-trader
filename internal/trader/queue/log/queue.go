package log

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-trader/internal/trader/queue"
	"time"
)

type marketQueue struct {
	logger *logrus.Logger
}

func (q *marketQueue) ReceiveMarkets() <-chan *queue.EventMarket {
	ch := make(chan *queue.EventMarket, 100)

	q.logger.Infof("Pretending to poll for messages from queue...")

	go q.simulate(ch)

	return ch
}

func (q *marketQueue) simulate(ch chan<- *queue.EventMarket) {
	time.Sleep(10 * time.Second)

	q.logger.Infof("..polling complete.")

	close(ch)
}

func NewMarketQueue(l *logrus.Logger) queue.MarketQueue {
	return &marketQueue{logger: l}
}
