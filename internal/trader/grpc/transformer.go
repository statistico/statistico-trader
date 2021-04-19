package grpc

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Transform a strategy.Trade struct into a statistico StrategyTrade struct to be consumed within the grpc package
func transformStrategyTrade(t *strategy.Trade) *statistico.StrategyTrade {
	return &statistico.StrategyTrade{
		MarketName:    t.MarketName,
		RunnerName:    t.RunnerName,
		RunnerPrice:   t.Price,
		Side:          statistico.SideEnum(statistico.SideEnum_value[t.Side]),
		EventId:       t.EventID,
		CompetitionId: t.CompetitionID,
		SeasonId:      t.SeasonID,
		EventDate:     timestamppb.New(t.EventDate),
		Result:        statistico.TradeResultEnum(statistico.TradeResultEnum_value[string(t.Result)]),
	}
}

// Transform statistico ResultFilter structs into trade ResultFilter structs to be consumed within the trade package
func transformResultFilters(f []*statistico.ResultFilter) []*strategy.ResultFilter {
	r := []*strategy.ResultFilter{}

	for _, filter := range f {
		s := &strategy.ResultFilter{
			Team:   filter.Team.String(),
			Result: filter.Result.String(),
			Games:  uint8(filter.Games),
			Venue:  filter.Venue.String(),
		}

		r = append(r, s)
	}

	return r
}

// Transform statistico StatFilter structs into trade StatFilter structs to be consumed within the trade package
func transformStatFilters(f []*statistico.StatFilter) []*strategy.StatFilter {
	s := []*strategy.StatFilter{}

	for _, filter := range f {
		sf := &strategy.StatFilter{
			Stat:    filter.Stat.String(),
			Team:    filter.Team.String(),
			Action:  filter.Action.String(),
			Games:   uint8(filter.Games),
			Measure: filter.Measure.String(),
			Metric:  filter.Metric.String(),
			Value:   filter.Value,
			Venue:   filter.Venue.String(),
		}

		s = append(s, sf)
	}

	return s
}
