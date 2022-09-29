package mock

import (
	"context"
	"time"

	"github.com/quarksgroup/sms-client/token"
)

type TokenSource struct {
	token *token.Token
}

func NewMockTokenSource() *TokenSource {
	return &TokenSource{
		token: &token.Token{
			Token:   "token",
			Refresh: "refresh",
			Expires: time.Now().Add(time.Hour),
		},
	}
}

func (t *TokenSource) Token(ctx context.Context) (*token.Token, error) {
	return t.token, nil
}

// check if TokenSource implements token.TokenSource
var _ token.TokenSource = (*TokenSource)(nil)
