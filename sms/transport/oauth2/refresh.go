package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/quarksgroup/sms-client/sms"
)

// expiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late
// expirations due to client-server time mismatches.
const expiryDelta = time.Minute

// Refresher is an http.RoundTripper that refreshes oauth
// tokens, wrapping a base RoundTripper and refreshing the
// token if expired.
type Refresher struct {
	Endpoint string

	Source sms.TokenSource
	Client *http.Client
}

// Token returns a token. If the token is missing or
// expired, the token is refreshed.
func (t *Refresher) Token(ctx context.Context) (*sms.Token, error) {
	token, err := t.Source.Token(ctx)
	if err != nil {
		return nil, err
	}
	if !expired(token) {
		return token, nil
	}
	err = t.Refresh(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Refresh refreshes the expired token.
func (t *Refresher) Refresh(token *sms.Token) error {
	in := tokenRefresh{
		Refresh: token.Refresh,
	}

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", t.Endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}
	res, err := t.client().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		out := new(tokenError)
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return err
		}
		return out
	}

	out := new(tokenGrant)
	err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return err
	}

	token.Token = out.Access
	token.Refresh = out.Refresh
	token.Expires = out.Expires

	return nil
}

// client returns the http transport. If no base client
// is configured, the default client is returned.
func (t *Refresher) client() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return http.DefaultClient
}

// expired reports whether the token is expired.
func expired(token *sms.Token) bool {
	if len(token.Refresh) == 0 {
		return false
	}
	if token.Expires.IsZero() && len(token.Token) != 0 {
		return false
	}
	return token.Expires.Add(-expiryDelta).
		Before(time.Now())
}

// tokenGrant is the token returned by the token endpoint.
type tokenGrant struct {
	Success bool      `json:"success"`
	Access  string    `json:"access_token"`
	Refresh string    `json:"refresh_token"`
	Expires time.Time `json:"expires_at"`
}

type tokenRefresh struct {
	Refresh string `json:"refresh_token"`
}

type tokenError struct {
	Message string `json:"message"`
}

func (t *tokenError) Error() string {
	return t.Message
}
