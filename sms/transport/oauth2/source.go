package oauth2

import (
	"context"

	"github.com/quarksgroup/sms-client/sms"
)

// StaticTokenSource returns a TokenSource that always
// returns the same token. Because the provided token t
// is never refreshed, StaticTokenSource is only useful
// for tokens that never expire.
func StaticTokenSource(t *sms.Token) sms.TokenSource {
	return staticTokenSource{t}
}

type staticTokenSource struct {
	token *sms.Token
}

func (s staticTokenSource) Token(context.Context) (*sms.Token, error) {
	return s.token, nil
}

// ContextTokenSource returns a TokenSource that returns
// a token from the http.Request context.
func ContextTokenSource() sms.TokenSource {
	return contextTokenSource{}
}

type contextTokenSource struct {
}

func (s contextTokenSource) Token(ctx context.Context) (*sms.Token, error) {
	token, _ := ctx.Value(sms.TokenKey{}).(*sms.Token)
	return token, nil
}
