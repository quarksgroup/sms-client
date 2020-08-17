package oauth2

import (
	"context"
	"testing"
	"time"

	"github.com/quarksgroup/sms-client/sms"
	"gopkg.in/h2non/gock.v1"
)

func TestRefresh(t *testing.T) {
	gock.New("https://messaging.fdibiz.com/api/v1/auth").
		Post("/refresh").
		Reply(200).
		File("testdata/token.json")

	before := &sms.Token{
		Refresh: "3a2bfce4cb9b0f",
	}

	r := Refresher{
		Endpoint: "https://messaging.fdibiz.com/api/v1/auth/refresh",
		Source:   StaticTokenSource(before),
	}

	ctx := context.Background()
	after, err := r.Token(ctx)
	if err != nil {
		t.Error(err)
	}

	if after.Token != "9698fa6a8113b3" {
		t.Errorf("Expect access token updated")
	}
	if after.Expires.IsZero() {
		t.Errorf("Expect access token expiry updated")
	}
	if after.Refresh != "3a2bfce4cb9b0f" {
		t.Errorf("Expect refresh token not changed, got %s", after.Refresh)
	}
}

func TestRefresh_NotExpired(t *testing.T) {
	before := &sms.Token{
		Token: "6084984dab20e6",
	}
	r := Refresher{
		Endpoint: "https://messaging.fdibiz.com/api/v1/auth/refresh",
		Source:   StaticTokenSource(before),
	}

	ctx := context.Background()
	after, err := r.Token(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if after == nil {
		t.Errorf("Expected Token, got nil")
		return
	}
	if after.Token != "6084984dab20e6" {
		t.Errorf("Expect Token not refreshed")
	}
}

func TestExpired(t *testing.T) {
	tests := []struct {
		token   *sms.Token
		expired bool
	}{
		{
			expired: false,
			token: &sms.Token{
				Token:   "12345",
				Refresh: "",
			},
		},
		{
			expired: false,
			token: &sms.Token{
				Token:   "12345",
				Refresh: "",
				Expires: time.Now().Add(-time.Hour),
			},
		},
		{
			expired: false,
			token: &sms.Token{
				Token:   "12345",
				Refresh: "54321",
			},
		},
		{
			expired: false,
			token: &sms.Token{
				Token:   "12345",
				Refresh: "54321",
				Expires: time.Now().Add(time.Hour),
			},
		},
		// missing access token
		{
			expired: true,
			token: &sms.Token{
				Token:   "",
				Refresh: "54321",
			},
		},
		// token expired
		{
			expired: true,
			token: &sms.Token{
				Token:   "12345",
				Refresh: "54321",
				Expires: time.Now().Add(-time.Second),
			},
		},
		// this token is not expired, however, it is within
		// the default 1 minute expiry window.
		{
			expired: true,
			token: &sms.Token{
				Token:   "12345",
				Refresh: "54321",
				Expires: time.Now().Add(time.Second * 30),
			},
		},
	}

	for i, test := range tests {
		if got, want := expired(test.token), test.expired; got != want {
			t.Errorf("Want token expired %v, got %v at index %d", want, got, i)
		}
	}
}
