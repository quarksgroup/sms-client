package sms

import (
	"context"
)

type (
	// Balance is the number of available message credit
	Balance struct {
		//The number of credits on your account
		Actual int64

		//The number of credits available for sending messages taking into account any pending holds
		Available int64

		// The date for which the balance was queried(Format:2019-01-01)
		Date string
	}
	// BalanceService ...
	BalanceService interface {
		// Current returns the current balance (as of now)
		Current(context.Context) (*Balance, *Response, error)

		// At returns the balance at a given date (eg;2020-09-02).
		At(context.Context, string) (*Balance, *Response, error)
	}
)
