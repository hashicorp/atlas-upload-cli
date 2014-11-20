package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	atlas "github.com/hashicorp/atlas-go"
)

func TestUpload_pending(t *testing.T) {
	t.Skip("not ready yet")
}

func TestAtlasClient_noURL(t *testing.T) {
	client, err := atlasClient(&UploadOpts{})
	if err != nil {
		t.Fatal(err)
	}

	expected := atlas.DefaultClient()
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

	expected, err := atlas.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
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

func TestProcess_errCh(t *testing.T) {
	doneCh, errCh := make(chan struct{}), make(chan error)
	go process(func() error {
		return fmt.Errorf("catastrophic failure")
	}, doneCh, errCh)

	select {
	case <-doneCh:
		t.Fatal("did not expect doneCh to receive data")
	case <-errCh:
		break
	case <-time.After(1 * time.Second):
		t.Fatal("no data returned in 1 second")
	}
}

func TestProcess_doneCh(t *testing.T) {
	doneCh, errCh := make(chan struct{}), make(chan error)
	go process(func() error {
		return nil
	}, doneCh, errCh)

	select {
	case <-doneCh:
		break
	case err := <-errCh:
		t.Fatal(err)
	case <-time.After(1 * time.Second):
		t.Fatal("no data returned in 1 second")
	}
}
