package fdi

import (
	"context"
	"fmt"

	"github.com/quarksgroup/sms-client/client"
)

// Balance ...
type Balance struct {
	//The number of credits on your account
	Actual int64

	//The number of credits available for sending messages taking into account any pending holds
	Available int64

	// The date for which the balance was queried(Format:2019-01-01)
	Date string
}

func (c *Client) Balance(ctx context.Context) (*Balance, *client.Response, error) {

	endpoint := "balance/now"
	out := new(currentBalance)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertCurrent(out), res, err

}

func (c *Client) BalanceAt(ctx context.Context, date string) (*Balance, *client.Response, error) {
	endpoint := fmt.Sprintf("balance/%s/closing", date)
	out := new(balanceAt)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertAt(out), res, err

}

func convertCurrent(current *currentBalance) *Balance {
	return &Balance{
		Actual:    current.Balance.Actual,
		Available: current.Balance.Available,
	}
}

func convertAt(bal *balanceAt) *Balance {
	return &Balance{
		Actual: bal.Balance,
		Date:   bal.Date,
	}
}

type currentBalance struct {
	Success bool    `json:"success"`
	Balance balance `json:"balance"`
}

type balanceAt struct {
	Success bool   `json:"success"`
	Balance int64  `json:"balance"`
	Date    string `json:"date"`
}

type balance struct {
	Actual    int64 `json:"available_balance"`
	Available int64 `json:"actual_balance"`
}
