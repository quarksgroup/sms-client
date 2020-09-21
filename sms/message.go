package sms

import "context"

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

	// SendService is the principal interface.
	// it sends a message returns a delivery report an serializable response
	// and if there is a problem returns an error
	SendService interface {
		Send(context.Context, Message) (*Report, *Response, error)
	}
)
