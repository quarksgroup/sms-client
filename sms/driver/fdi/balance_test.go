package fdi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/sms-client/sms"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestBalanceCurrent(t *testing.T) {
	defer gock.Off()

	gock.New("https://messaging.fdibiz.com/api/v1").
		Get("/balance/now").
		Reply(200).
		Type("application/json").
		File("testdata/balance.now.json")

	client := NewDefault()

	got, _, err := client.Balance.Current(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(sms.Balance)
	raw, _ := ioutil.ReadFile("testdata/balance.now.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestBalanceAt(t *testing.T) {
	defer gock.Off()

	var at = "2019-01-01"

	gock.New("https://messaging.fdibiz.com/api/v1").
		Get(fmt.Sprintf("/balance/%s/closing", at)).
		Reply(200).
		Type("application/json").
		File("testdata/balance.at.json")

	client := NewDefault()

	got, _, err := client.Balance.At(context.Background(), at)

	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(sms.Balance)
	raw, _ := ioutil.ReadFile("testdata/balance.at.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
