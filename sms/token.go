package sms

import (
	"context"
	"time"
)

type (
	// Token represents the credentials used to authorize
	// the requests to access protected resources.
	Token struct {
		Token   string
		Refresh string
		Expires time.Time
	}

	// TokenSource returns a token.
	TokenSource interface {
		Token(context.Context) (*Token, error)
	}

	// AuthService ..
	AuthService interface {
		// Login the underlying API and get an JWT token
		Login(ctx context.Context, id, secret string) (*Token, *Response, error)

		// Refresh the oauth2 token
		Refresh(ctx context.Context, token *Token, force bool) (*Token, *Response, error)
	}

	// TokenKey is the key to use with the context.WithValue
	// function to associate an Token value with a context.
	TokenKey struct{}
)

// WithContext returns a copy of parent in which the token value is set
func WithContext(parent context.Context, token *Token) context.Context {
	return context.WithValue(parent, TokenKey{}, token)
}

// TokenFrom returns the login token rom the context.
func TokenFrom(ctx context.Context) *Token {
	token, _ := ctx.Value(TokenKey{}).(*Token)
	return token
}
