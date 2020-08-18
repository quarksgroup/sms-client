package fdi

import (
	"context"
	"time"

	"github.com/quarksgroup/sms-client/sms"
)

var _ (sms.LoginService) = (*loginService)(nil)

type loginService struct {
	client *wrapper
}

func (s *loginService) Login(ctx context.Context, id, secret string) (*sms.Token, *sms.Response, error) {
	endpoint := "auth"
	in := login{
		User:     id,
		Password: secret,
	}
	out := new(token)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertToken(out), res, err
}

type login struct {
	User     string `json:"api_username"`
	Password string `json:"api_password"`
}

type token struct {
	Success bool      `json:"success"`
	Access  string    `json:"access_token"`
	Refresh string    `json:"refresh_token"`
	Expires time.Time `json:"expires_at"`
}

func convertToken(from *token) *sms.Token {
	return &sms.Token{
		Token:   from.Access,
		Refresh: from.Refresh,
		Expires: from.Expires,
	}
}
