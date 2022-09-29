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

func TestStatCurrent(t *testing.T) {
	defer gock.Off()

	gock.New("https://messaging.fdibiz.com/api/v1").
		Get("/stats/now").
		Reply(200).
		Type("application/json").
		File("testdata/stats.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
	}

	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.Current(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))

	want := new(Stats)

	raw, _ := ioutil.ReadFile("testdata/stats.json.golden")

	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestStatsAt(t *testing.T) {
	defer gock.Off()

	var at = "2019-01-01"

	gock.New("https://messaging.fdibiz.com/api/v1").
		Get(fmt.Sprintf("/stats/%s", at)).
		Reply(200).
		Type("application/json").
		File("testdata/stats.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
	}

	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.At(context.Background(), at)

	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))

	want := new(Stats)

	raw, _ := ioutil.ReadFile("testdata/stats.json.golden")

	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
