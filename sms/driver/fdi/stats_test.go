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

func TestStatCurrent(t *testing.T) {
	defer gock.Off()

	gock.New("https://messaging.fdibiz.com/api/v1").
		Get("/stats/now").
		Reply(200).
		Type("application/json").
		File("testdata/stats.json")

	client := NewDefault()

	got, _, err := client.Stats.Current(context.Background())
	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(sms.Stats)
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

	client := NewDefault()

	got, _, err := client.Stats.At(context.Background(), at)
	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(sms.Stats)
	raw, _ := ioutil.ReadFile("testdata/stats.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
