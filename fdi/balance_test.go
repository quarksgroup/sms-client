package fdi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/sms-client/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

// TestBalanceCurrent tests the balance at the current time
func TestBalanceCurrent(t *testing.T) {
	defer gock.Off()

	gock.New("https://messaging.fdibiz.com/api/v1").
		Get("/balance/now").
		Reply(200).
		Type("application/json").
		File("testdata/balance.now.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
	}

	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.Balance(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(Balance)
	raw, _ := ioutil.ReadFile("testdata/balance.now.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

// TestBalanceAt tests the balance at a given time
func TestBalanceAt(t *testing.T) {
	defer gock.Off()

	var at = "2019-01-01"

	gock.New("https://messaging.fdibiz.com/api/v1").
		Get(fmt.Sprintf("/balance/%s/closing", at)).
		Reply(200).
		Type("application/json").
		File("testdata/balance.at.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
	}

	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.BalanceAt(context.Background(), at)

	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(Balance)
	raw, _ := ioutil.ReadFile("testdata/balance.at.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
