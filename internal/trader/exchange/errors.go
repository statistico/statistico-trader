package exchange

import "fmt"

type ClientError struct {
	Action string
	Err    error
}

func (c *ClientError) Error() string {
	return fmt.Sprintf("error making '%s' request: %s", c.Action, c.Err.Error())
}

type OrderFailureError struct {
	MarketID  string
	RunnerID  uint64
	Status    string
	ErrorCode string
}

func (o *OrderFailureError) Error() string {
	return fmt.Sprintf(
		"error placing order for market '%s' and runner '%d'. Code: '%s' and Status: '%s'",
		o.MarketID,
		o.RunnerID,
		o.ErrorCode,
		o.Status,
	)
}

type InvalidResponseError struct {
	Message string
}

func (i *InvalidResponseError) Error() string {
	return fmt.Sprintf("invalid response in exchange client: %s", i.Message)
}

type UnmatchedError struct {
	MarketID  string
	RunnerID  uint64
	Status    string
}

func (u *UnmatchedError) Error() string {
	return fmt.Sprintf("trade unmatched for market '%s' and runner '%d' with status '%s'", u.MarketID, u.RunnerID, u.Status)
}
