package fdi

import (
	"context"
	"time"

	"github.com/quarksgroup/sms-client/token"
)

// tokenSource implements token.TokenSource
type tokenSource struct {
	token  *token.Token
	client *Client
	id     string
	secret string
}

// newTokenSource creates a new tokenSource instance backed by the  http.Client instance and the given credentials
func newTokenSource(client *Client, cfg *Config) (token.TokenSource, error) {
	tks := &tokenSource{
		client: client,
		id:     cfg.ClientId,
		secret: cfg.Secret,
	}
	token, err := tks.Login(context.Background(), cfg.ClientId, cfg.Secret)
	if err != nil {
		return nil, err
	}
	tks.token = token
	return tks, nil
}

// Token returns the current token or refreshes it if it has expired
func (tk *tokenSource) Token(ctx context.Context) (*token.Token, error) {
	if tk.token != nil {
		if tk.token.Expires.Before(time.Now().Local()) {
			token, err := tk.Refresh(ctx, tk.token)
			if err != nil {
				return nil, err
			}
			tk.token = token
		}
		return tk.token, nil
	}
	return tk.Login(ctx, tk.id, tk.secret)
}

// Login returns a new token from the given credentials or an error if the credentials are invalid or other sms api error happen
func (tk *tokenSource) Login(ctx context.Context, id, secret string) (*token.Token, error) {
	endpoint := "auth"
	in := login{
		User:     id,
		Password: secret,
	}
	out := new(tokenGrant)
	_, err := tk.client.do(ctx, "POST", endpoint, in, out, false)
	return convertToken(out), err
}

// Refresh returns a new token from the given refresh token or an error if the refresh token is invalid or other sms api error happen
func (tk *tokenSource) Refresh(ctx context.Context, token *token.Token) (*token.Token, error) {
	endpoint := "auth/refresh"
	in := tokenRefresh{
		Refresh: token.Refresh,
	}
	out := new(tokenGrant)
	_, err := tk.client.do(ctx, "POST", endpoint, in, out, false)
	return convertToken(out), err
}

// Convert tokenGrant to token.Token
func convertToken(tk *tokenGrant) *token.Token {
	return &token.Token{
		Token:   tk.Access,
		Refresh: tk.Refresh,
		Expires: tk.Expires,
	}
}

// login credentials request
type login struct {
	User     string `json:"api_username"`
	Password string `json:"api_password"`
}

// tokenGrant response
type tokenGrant struct {
	Success bool      `json:"success"`
	Access  string    `json:"access_token"`
	Refresh string    `json:"refresh_token"`
	Expires time.Time `json:"expires_at"`
}

// refresh token request
type tokenRefresh struct {
	Refresh string `json:"refresh_token"`
}
