package fdi

import (
	"context"
	"fmt"
	"time"

	"github.com/quarksgroup/sms-client/client"
)

// Stats contains information on message deliverry
type Stats struct {
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
func (c *Client) At(ctx context.Context, at string) (*Stats, *client.Response, error) {
	endpoint := fmt.Sprintf("stats/%s", at)
	out := new(stats)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertStats(out), res, err
}

func (c *Client) Current(ctx context.Context) (*Stats, *client.Response, error) {
	endpoint := "stats/now"
	out := new(stats)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertStats(out), res, err
}

type stats struct {
	Success bool      `json:"success"`
	Date    time.Time `json:"date"`
	MT      outBox    `json:"mt"`
	MO      inBox     `json:"mo"`
}

type outBox struct {
	Sent      int64 `json:"sent"`
	Delivered int64 `json:"delivered"`
	Failed    int64 `json:"failed"`
	Pending   int64 `json:"pending"`
}

type inBox struct {
	Recieved int64 `json:"received"`
}

func convertStats(from *stats) *Stats {
	return &Stats{
		Sent:      from.MT.Sent,
		Delivered: from.MT.Delivered,
		Pending:   from.MT.Pending,
		Failed:    from.MT.Failed,
		Received:  from.MO.Recieved,
		Date:      from.Date,
	}
}
