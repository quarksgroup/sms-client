package fdi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/quarksgroup/sms-client/sms"
)

// New returns a new FDI sms API client.
func New(uri string) (*sms.Client, error) {

	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	client := &wrapper{new(sms.Client)}
	client.BaseURL = base
	// initialize services
	client.Driver = sms.DriverFdi
	client.Balance = &balanceService{client}
	client.Message = &sendService{client}
	client.Auth = &loginService{client}
	client.Stats = &statsService{client}

	return client.Client, nil
}

// NewDefault returns a new FDI API client using the
// default "https://messaging.fdibiz.com/api/v1" address.
func NewDefault() *sms.Client {
	client, _ := New("https://messaging.fdibiz.com/api/v1")
	return client
}

type wrapper struct {
	*sms.Client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*sms.Response, error) {
	req := &sms.Request{
		Method: method,
		Path:   path,
		// Header: map[string][]string{"Content-Type": {"application/json"}},
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(in)
		req.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		req.Body = buf
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status > 299 {
		err := new(Error)
		err.Code = res.Status
		_ = json.NewDecoder(res.Body).Decode(err)
		return res, err
	}

	if out == nil {
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// Error represents a Github error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
