package exchange

import "fmt"

type ClientError struct {
	action string
	err    error
}

func (c *ClientError) Error() string {
	return fmt.Sprintf("error making '%s' request: %s", c.action, c.err.Error())
}

type OrderFailureError struct {
	marketID  string
	runnerID  uint64
	status    string
	errorCode string
	stake     float32
	price     float32
}

func (o *OrderFailureError) Error() string {
	return fmt.Sprintf(
		"error placing order for market %s and runner %d. Code: %s and Status: %s and Stake: %.2f and Price: %.2f",
		o.marketID,
		o.runnerID,
		o.errorCode,
		o.status,
		o.stake,
		o.price,
	)
}

type InvalidExchangeError struct {
	exchange string
}

func (i *InvalidExchangeError) Error() string {
	return fmt.Sprintf("exchange '%s' is not supported", i.exchange)
}

type InvalidResponseError struct {
	message string
}

func (i *InvalidResponseError) Error() string {
	return fmt.Sprintf("invalid response in exchange client: %s", i.message)
}

type UnmatchedError struct {
	marketID  string
	runnerID  uint64
	status    string
}

func (u *UnmatchedError) Error() string {
	return fmt.Sprintf("trade unmatched for market '%s' and runner '%d' with status '%s'", u.marketID, u.runnerID, u.status)
}
