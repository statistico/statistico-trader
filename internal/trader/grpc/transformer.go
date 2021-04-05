package grpc

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
)

// Transform a market Trade struct into a statistico StrategyTrade struct to be consumed within the grpc package
//func transformTradeResultToStrategyTrade(t *market.Trade) (*statistico.StrategyTrade, error) {
//	date, _ := ptypes.TimestampProto(t.EventDate)
//
//	if t.Result == nil {
//		return nil, fmt.Errorf("event %d does not contain a result", t.EventId)
//	}
//
//	return &statistico.StrategyTrade{
//		MarketName:    t.MarketName,
//		RunnerName:    t.RunnerName,
//		RunnerPrice:   t.RunnerPrice,
//		Side:          statistico.SideEnum(statistico.SideEnum_value[t.Side]),
//		EventId:       t.EventId,
//		CompetitionId: t.CompetitionId,
//		SeasonId:      t.SeasonId,
//		EventDate:     date,
//		Result:        statistico.TradeResultEnum(statistico.TradeResultEnum_value[*t.Result]),
//	}, nil
//}

// Transform statistico ResultFilter structs into trade ResultFilter structs to be consumed within the trade package
func transformResultFilters(f []*statistico.ResultFilter) []*trader.ResultFilter {
	r := []*trader.ResultFilter{}

	for _, filter := range f {
		s := &trader.ResultFilter{
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
func transformStatFilters(f []*statistico.StatFilter) []*trader.StatFilter {
	s := []*trader.StatFilter{}

	for _, filter := range f {
		sf := &trader.StatFilter{
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
