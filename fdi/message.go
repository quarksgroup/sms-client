package fdi

import (
	"context"
	"fmt"

	"github.com/quarksgroup/sms-client/client"
)

type (
	// Message contains all the details about the sms to be sent
	Message struct {
		ID         string // Unique id for the current message
		Body       string // The content of the message
		Report     string // Callback url to report back to(Optional)
		Sender     string
		Recipients []string // The recipients of this particular message
	}

	// Report back sent message details
	Report struct {
		ID   string
		Cost int64
	}
)

// it sends a message returns a delivery report an serializable response
// and if there is a problem returns an error
func (c *Client) Send(ctx context.Context, message Message) (*Report, *client.Response, error) {
	n := len(message.Recipients)
	if n < 1 {
		return nil, nil, fmt.Errorf("can't send message to zero recipients")
	}
	if n == 1 {
		return c.Single(ctx, message)
	}
	return c.Bulk(ctx, message)
}

func (c *Client) Single(ctx context.Context, message Message) (*Report, *client.Response, error) {
	endpoint := "mt/single"
	in := singleSend{
		Reference: message.ID,
		Body:      message.Body,
		Sender:    message.Sender,
		Dlr:       message.Report,
		MSISDN:    message.Recipients[0],
	}
	out := new(reportSingle)
	res, err := c.do(ctx, "POST", endpoint, in, out, true)
	return convertSingle(out), res, err
}

func (c *Client) Bulk(ctx context.Context, message Message) (*Report, *client.Response, error) {
	endpoint := "mt/bulk"
	in := bulkSend{
		Reference: message.ID,
		Body:      message.Body,
		Sender:    message.Sender,
		Dlr:       message.Report,
		MSISDN:    message.Recipients,
	}
	out := new(reportBulk)
	res, err := c.do(ctx, "POST", endpoint, in, out, true)
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

func convertSingle(from *reportSingle) *Report {
	return &Report{
		ID:   from.MessageReference,
		Cost: from.Cost,
	}
}

func convertBulk(from *reportBulk) *Report {
	return &Report{
		ID:   from.MessageReference,
		Cost: from.Cost,
	}
}
