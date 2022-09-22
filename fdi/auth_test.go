package fdi

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"testing"
// 	"time"

// 	"github.com/google/go-cmp/cmp"
// 	"github.com/quarksgroup/sms-client/sms"
// 	"github.com/quarksgroup/sms-client/token"
// 	"github.com/stretchr/testify/require"
// 	"gopkg.in/h2non/gock.v1"
// )

// func TestLogin(t *testing.T) {
// 	defer gock.Off()

// 	gock.New("https://messaging.fdibiz.com/api/v1").
// 		Post("/auth").
// 		Reply(200).
// 		Type("application/json").
// 		File("testdata/token.json")
// 	client, err := NewDefault("xxxx", "kdskslas")

// 	got, _, err := client.Auth.Login(context.Background(), "id", "secret")

// 	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

// 	want := new(token.Token)
// 	raw, _ := ioutil.ReadFile("testdata/token.json.golden")
// 	_ = json.Unmarshal(raw, want)

// 	if diff := cmp.Diff(got, want); diff != "" {
// 		t.Errorf("Unexpected Results")
// 		t.Log(diff)
// 	}
// }

// func TestRefresh(t *testing.T) {
// 	gock.New("https://messaging.fdibiz.com/api/v1/auth").
// 		Post("/refresh").
// 		Reply(200).
// 		File("testdata/token.json")

// 	client := NewDefault()

// 	expired := &sms.Token{
// 		Refresh: "3a2bfce4cb9b0f",
// 	}

// 	after, _, err := client.Auth.Refresh(context.Background(), expired, false)
// 	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

// 	want := new(sms.Token)
// 	raw, _ := ioutil.ReadFile("testdata/token.json.golden")
// 	_ = json.Unmarshal(raw, want)

// 	if after.Token != want.Token {
// 		t.Errorf("Expect access token updated")
// 	}
// 	if after.Expires.IsZero() {
// 		t.Errorf("Expect access token expiry updated")
// 	}
// 	if after.Refresh != want.Refresh {
// 		t.Errorf("Expect refresh token not changed, got %s", after.Refresh)
// 	}

// }

// func TestRefresh_NotExpired(t *testing.T) {
// 	client := NewDefault()

// 	before := &sms.Token{
// 		Token: "6084984dab20e6",
// 	}

// 	after, _, err := client.Auth.Refresh(context.Background(), before, false)
// 	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

// 	if after == nil {
// 		t.Errorf("Expected Token, got nil")
// 		return
// 	}
// 	if after.Token != "6084984dab20e6" {
// 		t.Errorf("Expect Token not refreshed")
// 	}
// }

// func TestExpired(t *testing.T) {
// 	tests := []struct {
// 		token   *sms.Token
// 		expired bool
// 	}{
// 		{
// 			expired: false,
// 			token: &sms.Token{
// 				Token:   "12345",
// 				Refresh: "",
// 			},
// 		},
// 		{
// 			expired: false,
// 			token: &sms.Token{
// 				Token:   "12345",
// 				Refresh: "",
// 				Expires: time.Now().Add(-time.Hour),
// 			},
// 		},
// 		{
// 			expired: false,
// 			token: &sms.Token{
// 				Token:   "12345",
// 				Refresh: "54321",
// 			},
// 		},
// 		{
// 			expired: false,
// 			token: &sms.Token{
// 				Token:   "12345",
// 				Refresh: "54321",
// 				Expires: time.Now().Add(time.Hour),
// 			},
// 		},
// 		// missing access token
// 		{
// 			expired: true,
// 			token: &sms.Token{
// 				Token:   "",
// 				Refresh: "54321",
// 			},
// 		},
// 		// token expired
// 		{
// 			expired: true,
// 			token: &sms.Token{
// 				Token:   "12345",
// 				Refresh: "54321",
// 				Expires: time.Now().Add(-time.Second),
// 			},
// 		},
// 		// this token is not expired, however, it is within
// 		// the default 1 minute expiry window.
// 		{
// 			expired: true,
// 			token: &sms.Token{
// 				Token:   "12345",
// 				Refresh: "54321",
// 				Expires: time.Now().Add(time.Second * 30),
// 			},
// 		},
// 	}

// 	for i, test := range tests {
// 		if got, want := expired(test.token), test.expired; got != want {
// 			t.Errorf("Want token expired %v, got %v at index %d", want, got, i)
// 		}
// 	}
// }
