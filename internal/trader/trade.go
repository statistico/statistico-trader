package trader

type TradeWriter interface {
	Insert(t *Trade) error
}
