package fdi

import (
	"context"
	"fmt"

	"github.com/quarksgroup/sms-client/sms"
)

var _ (sms.BalanceService) = (*balanceService)(nil)

type balanceService struct {
	client *wrapper
}

func (s *balanceService) Current(ctx context.Context) (*sms.Balance, *sms.Response, error) {
	endpoint := "balance/now"
	out := new(currentBalance)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return convertCurrent(out), res, err
}

func (s *balanceService) At(ctx context.Context, at string) (*sms.Balance, *sms.Response, error) {
	endpoint := fmt.Sprintf("balance/%s/closing", at)
	out := new(balanceAt)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return convertAt(out), res, err
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

func convertCurrent(from *currentBalance) *sms.Balance {
	return &sms.Balance{
		Actual:    from.Balance.Actual,
		Available: from.Balance.Available,
	}
}

func convertAt(from *balanceAt) *sms.Balance {
	return &sms.Balance{
		Actual: from.Balance,
		Date:   from.Date,
	}
}
