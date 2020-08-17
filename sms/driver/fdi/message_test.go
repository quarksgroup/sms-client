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

func TestSingleSend(t *testing.T) {
	gock.New("https://messaging.fdibiz.com/api/v1").
		Post("/mt/single").
		Reply(200).
		Type("application/json").
		File("testdata/message.single.json")

	client := NewDefault()

	message := sms.Message{
		ID:         "30bb083a-ae95-43b9-8ed5-051693d018af",
		Body:       "Hello world",
		Sender:     "Paypack",
		Recipients: []string{"0789640100"},
	}

	got, _, err := client.Message.Send(context.Background(), message)
	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(sms.Report)
	raw, _ := ioutil.ReadFile("testdata/message.single.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestBulkSend(t *testing.T) {
	gock.New("https://messaging.fdibiz.com/api/v1").
		Post("/mt/bulk").
		Reply(200).
		Type("application/json").
		File("testdata/message.bulk.json")

	client := NewDefault()

	message := sms.Message{
		ID:         "30bb083a-ae95-43b9-8ed5-051693d018af",
		Body:       "Hello world",
		Sender:     "Paypack",
		Recipients: []string{"0789640100", "0783205104"},
	}

	got, _, err := client.Message.Send(context.Background(), message)
	require.Nil(t, err, fmt.Sprintf("unexpected error '%v'", err))
	want := new(sms.Report)
	raw, _ := ioutil.ReadFile("testdata/message.bulk.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
