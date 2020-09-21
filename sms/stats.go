package sms

import (
	"context"
	"time"
)

type (
	// Stats contains information on message deliverry
	Stats struct {
		// The total number of messages sent
		Sent int64

		// The number of messages successfully delivered
		Delivered int64

		// The number of messages that failed to be delivered
		Failed int64

		// The number of messages pending delivery
		Pending int64

		// The total number of messages received
		Received int64

		// The date at which this data is obtained
		Date time.Time
	}
	// StatsService handles information(stats) about your sms sending habits.
	StatsService interface {
		// At gets credit balance at any given date (eg;2020-09-02).
		At(ctx context.Context, at string) (*Stats, *Response, error)

		// Current gets the current credit Balance.
		Current(ctx context.Context) (*Stats, *Response, error)
	}
)
