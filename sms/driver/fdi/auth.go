package fdi

import (
	"context"
	"time"

	"github.com/quarksgroup/sms-client/sms"
)

// expiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late
// expirations due to client-server time mismatches.
const expiryDelta = time.Minute

var _ (sms.AuthService) = (*loginService)(nil)

type loginService struct {
	client *wrapper
}

func (s *loginService) Login(ctx context.Context, id, secret string) (*sms.Token, *sms.Response, error) {
	endpoint := "auth"
	in := login{
		User:     id,
		Password: secret,
	}
	out := new(tokenGrant)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertToken(out), res, err
}

func (s *loginService) Refresh(ctx context.Context, token *sms.Token, force bool) (*sms.Token, *sms.Response, error) {
	endpoint := "auth/refresh"
	in := tokenRefresh{}
	if !expired(token) {
		return token, nil, nil
	}
	out := new(tokenGrant)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertToken(out), res, err
}

type login struct {
	User     string `json:"api_username"`
	Password string `json:"api_password"`
}

type tokenGrant struct {
	Success bool      `json:"success"`
	Access  string    `json:"access_token"`
	Refresh string    `json:"refresh_token"`
	Expires time.Time `json:"expires_at"`
}
type tokenRefresh struct {
	Refresh string `json:"refresh_token"`
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

func convertToken(from *tokenGrant) *sms.Token {
	return &sms.Token{
		Token:   from.Access,
		Refresh: from.Refresh,
		Expires: from.Expires,
	}
}
