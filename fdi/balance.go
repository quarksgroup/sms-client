package fdi

import (
	"context"
	"fmt"

	"github.com/quarksgroup/sms-client/client"
)

// Balance represents the balance of the account
type Balance struct {
	//The number of credits on your account
	Actual int64

	//The number of credits available for sending messages taking into account any pending holds
	Available int64

	// The date for which the balance was queried(Format:2019-01-01)
	Date string
}

// BalanceCurrent returns the current balance of the account
func (c *Client) Balance(ctx context.Context) (*Balance, *client.Response, error) {

	endpoint := "balance/now"
	out := new(currentBalance)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertCurrent(out), res, err

}

// BalanceAt returns the balance of the account at the given date
func (c *Client) BalanceAt(ctx context.Context, date string) (*Balance, *client.Response, error) {
	endpoint := fmt.Sprintf("balance/%s/closing", date)
	out := new(balanceAt)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertAt(out), res, err

}

// convertCurrent converts the currentBalance to Balance for our end use
func convertCurrent(current *currentBalance) *Balance {
	return &Balance{
		Actual:    current.Balance.Actual,
		Available: current.Balance.Available,
	}
}

// convertAt converts the balanceAt to Balance for our end use
func convertAt(bal *balanceAt) *Balance {
	return &Balance{
		Actual: bal.Balance,
		Date:   bal.Date,
	}
}

// currentBalance represents the current balance of the account for our end
type currentBalance struct {
	Success bool    `json:"success"`
	Balance balance `json:"balance"`
}

// balanceAt represents the balance of the account at the given date for our end
type balanceAt struct {
	Success bool   `json:"success"`
	Balance int64  `json:"balance"`
	Date    string `json:"date"`
}

// balance represents the balance of the account for sms gateway provider
type balance struct {
	Actual    int64 `json:"available_balance"`
	Available int64 `json:"actual_balance"`
}
