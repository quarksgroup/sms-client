package fdi

import (
	"context"
	"fmt"

	"github.com/quarksgroup/sms-client/sms"
)

var _ (sms.SendService) = (*sendService)(nil)

type sendService struct {
	client *wrapper
}

func (s *sendService) Send(ctx context.Context, message sms.Message) (*sms.Report, *sms.Response, error) {
	n := len(message.Recipients)
	if n < 1 {
		return nil, nil, fmt.Errorf("can't send message to zero recipients")
	}
	if n == 1 {
		return s.Single(ctx, message)
	}
	return s.Bulk(ctx, message)
}

func (s *sendService) Single(ctx context.Context, message sms.Message) (*sms.Report, *sms.Response, error) {
	endpoint := "mt/single"
	in := singleSend{
		Reference: message.ID,
		Body:      message.Body,
		Sender:    message.Sender,
		Dlr:       message.Report,
		MSISDN:    message.Recipients[0],
	}
	out := new(reportSingle)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertSingle(out), res, err
}

func (s *sendService) Bulk(ctx context.Context, message sms.Message) (*sms.Report, *sms.Response, error) {
	endpoint := "mt/bulk"
	in := bulkSend{
		Reference: message.ID,
		Body:      message.Body,
		Sender:    message.Sender,
		Dlr:       message.Report,
		MSISDN:    message.Recipients,
	}
	out := new(reportBulk)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertBulk(out), res, err
}

type singleSend struct {
	Reference string `json:"msgRef"`
	Body      string `json:"message"`
	MSISDN    string `json:"msisdn"`
	Dlr       string `json:"dlr,omitempty"`
	Sender    string `json:"sender_id,omitempty"`
}

type bulkSend struct {
	Reference string   `json:"msgRef"`
	Body      string   `json:"message"`
	MSISDN    []string `json:"msisdn_list"`
	Dlr       string   `json:"dlr,omitempty"`
	Sender    string   `json:"sender_id,omitempty"`
}

type reportSingle struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	Cost             int64  `json:"cost"`
	MessageReference string `json:"msgRef"`
	GatewayReference string `json:"gatewayRef"`
}

type reportBulk struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	Cost             int64  `json:"cost"`
	MessageReference string `json:"msgRef"`
	GatewayReference string `json:"gatewayRef"`
	Data             struct {
		Valid   int64 `json:"valid"`
		Invalid int64 `json:"invalid"`
	} `json:"data"`
}

func convertSingle(from *reportSingle) *sms.Report {
	return &sms.Report{
		ID:   from.MessageReference,
		Cost: from.Cost,
	}
}

func convertBulk(from *reportBulk) *sms.Report {
	return &sms.Report{
		ID:   from.MessageReference,
		Cost: from.Cost,
	}
}
