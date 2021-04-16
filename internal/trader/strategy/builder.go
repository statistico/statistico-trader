package strategy

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-odds-warehouse-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"sync"
)

type Builder interface {
	Build(ctx context.Context, q *BuilderQuery) <-chan *Trade
}

type builder struct {
	matcher    FilterMatcher
	parser     ResultParser
	oddsClient statisticooddswarehouse.MarketClient
	logger     *logrus.Logger
}

func (b *builder) Build(ctx context.Context, q *BuilderQuery) <-chan *Trade {
	ch := make(chan *Trade, 1000)

	go b.build(ctx, ch, q)

	return ch
}

func (b *builder) build(ctx context.Context, ch chan<- *Trade, q *BuilderQuery) {
	defer close(ch)
	var wg sync.WaitGroup

	req := buildMarketRequest(q)

	markets, errCh := b.oddsClient.MarketRunnerSearch(ctx, req)

	for mk := range markets {
		wg.Add(1)
		go b.handleMarket(ctx, ch, mk, q, &wg)
	}

	err := <- errCh

	if err != nil {
		b.logger.Errorf("error fetching market runners from odds warehouse: %s", err.Error())
	}

	wg.Wait()
}

func (b *builder) handleMarket(ctx context.Context, ch chan<- *Trade, mk *statistico.MarketRunner, q *BuilderQuery, wg *sync.WaitGroup) {
	query := MatcherQuery{
		EventID:       mk.EventId,
		ResultFilters: q.ResultFilters,
		StatFilters:   q.StatFilters,
	}

	matches, err := b.matcher.MatchesFilters(ctx, &query)

	if err != nil {
		b.logError(mk.MarketName, mk.RunnerName, mk.EventId, err)
		wg.Done()
		return
	}

	if matches {
		result, err := b.parser.Parse(ctx, mk.EventId, mk.MarketName, mk.RunnerName, q.Side)

		if err != nil {
			b.logError(mk.MarketName, mk.RunnerName, mk.EventId, err)
			wg.Done()
			return
		}

		tr := &Trade{
			MarketName:    mk.MarketName,
			RunnerName:    mk.RunnerName,
			EventID:       mk.EventId,
			CompetitionID: mk.CompetitionId,
			SeasonID:      mk.SeasonId,
			EventDate:     mk.EventDate.AsTime(),
			Exchange:      mk.Exchange,
			Price:         mk.Price.Value,
			Side:          q.Side,
			Result:        result,
		}

		ch <- tr
	}

	wg.Done()
}

func (b *builder) logError(market, runner string, eventID uint64, e error) {
	b.logger.Errorf(
		"error handling trade for market %s, runner %s and event %d: %+v",
		market,
		runner,
		eventID,
		e,
	)
}

func buildMarketRequest(q *BuilderQuery) *statistico.MarketRunnerRequest {
	req := statistico.MarketRunnerRequest{
		Market:         q.Market,
		Runner:         q.Runner,
		Line:           q.Line,
		Side:           statistico.SideEnum(statistico.SideEnum_value[q.Side]),
		CompetitionIds: q.CompetitionIDs,
		SeasonIds:      q.SeasonIDs,
	}

	if q.MinOdds != nil {
		req.MinOdds = &wrappers.FloatValue{Value: *q.MinOdds}
	}

	if q.MaxOdds != nil {
		req.MaxOdds = &wrappers.FloatValue{Value: *q.MaxOdds}
	}

	return &req
}

func NewBuilder(m FilterMatcher, p ResultParser, o statisticooddswarehouse.MarketClient, l *logrus.Logger) Builder {
	return &builder{
		matcher:    m,
		parser:     p,
		oddsClient: o,
		logger:     l,
	}
}
