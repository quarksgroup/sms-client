package sms

import "context"

type (
	// Message ...
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

	// SendService ...
	SendService interface {
		Send(context.Context, Message) (*Report, *Response, error)
	}
)
