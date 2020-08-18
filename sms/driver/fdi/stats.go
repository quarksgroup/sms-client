package fdi

import (
	"context"
	"fmt"
	"time"

	"github.com/quarksgroup/sms-client/sms"
)

type statsService struct {
	client *wrapper
}

func (s *statsService) At(ctx context.Context, at string) (*sms.Stats, *sms.Response, error) {
	endpoint := fmt.Sprintf("stats/%s", at)
	out := new(stats)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return convertStats(out), res, err
}

func (s *statsService) Current(ctx context.Context) (*sms.Stats, *sms.Response, error) {
	endpoint := "stats/now"
	out := new(stats)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
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

func convertStats(from *stats) *sms.Stats {
	return &sms.Stats{
		Sent:      from.MT.Sent,
		Delivered: from.MT.Delivered,
		Pending:   from.MT.Pending,
		Failed:    from.MT.Failed,
		Received:  from.MO.Recieved,
		Date:      from.Date,
	}
}
