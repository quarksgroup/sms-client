package oauth2

import (
	"context"
	"testing"

	"github.com/quarksgroup/sms-client/sms"
)

func TestContextTokenSource(t *testing.T) {
	source := ContextTokenSource()
	want := new(sms.Token)

	ctx := context.Background()
	ctx = context.WithValue(ctx, sms.TokenKey{}, want)
	got, err := source.Token(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if got != want {
		t.Errorf("Expect token retrieved from Context")
	}
}

func TestContextTokenSource_Nil(t *testing.T) {
	source := ContextTokenSource()

	ctx := context.Background()
	token, err := source.Token(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if token != nil {
		t.Errorf("Expect nil token from Context")
	}
}
