package fdi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/sms-client/token"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

//Response is the response from the sms gateway api for getting our token case
func TestLogin(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	t.Skip()

	gock.New("https://messaging.fdibiz.com/api/v1").
		Post("/auth").
		Reply(200).
		Type("application/json").
		File("testdata/token.json")
	client, err := NewDefault("client_id", "screte")

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, err := client.TokenSource.Token(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(token.Token)
	raw, _ := ioutil.ReadFile("testdata/token.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
