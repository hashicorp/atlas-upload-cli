package main

import (
	"reflect"
	"testing"

	"github.com/hashicorp/atlas-go/v1"
)

func TestUpload_pending(t *testing.T) {
	t.Skip("not ready yet")
}

func TestAtlasClient_noURL(t *testing.T) {
	client, err := atlasClient(&UploadOpts{})
	if err != nil {
		t.Fatal(err)
	}
	client.HTTPClient = nil

	expected := atlas.DefaultClient()
	expected.HTTPClient = nil
	if !reflect.DeepEqual(client, expected) {
		t.Fatalf("expected %+v to be %+v", client, expected)
	}
}

func TestAtlasClient_customURL(t *testing.T) {
	url := "https://atlas.company.com"
	client, err := atlasClient(&UploadOpts{URL: url})
	if err != nil {
		t.Fatal(err)
	}
	client.HTTPClient = nil

	expected, err := atlas.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	expected.HTTPClient = nil
	if !reflect.DeepEqual(client, expected) {
		t.Fatalf("expected %+v to be %+v", client, expected)
	}
}

func TestAtlasClient_token(t *testing.T) {
	token := "abcd1234"
	client, err := atlasClient(&UploadOpts{Token: token})
	if err != nil {
		t.Fatal(err)
	}

	if client.Token != token {
		t.Fatalf("expected %q to be %q", client.Token, token)
	}
}
