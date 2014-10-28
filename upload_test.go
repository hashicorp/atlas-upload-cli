package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	harmony "github.com/hashicorp/harmony-go"
)

func TestUpload_pending(t *testing.T) {
	t.Skip("not ready yet")
}

func TestParseSlug_emptyString(t *testing.T) {
	_, _, err := parseSlug("")
	if err == nil {
		t.Fatal("expected error, but nothing was returned")
	}

	expected := "missing slug"
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected %q to contain %q", err.Error(), expected)
	}
}

func TestParseSlug_noSlashes(t *testing.T) {
	_, _, err := parseSlug("bacon")
	if err == nil {
		t.Fatal("expected error, but nothing was returned")
	}

	expected := "malformed slug"
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected %q to contain %q", err.Error(), expected)
	}
}

func TestParseSlug_multipleSlashes(t *testing.T) {
	_, _, err := parseSlug("bacon/is/delicious/but/this/is/not/valid")
	if err == nil {
		t.Fatal("expected error, but nothing was returned")
	}

	expected := "malformed slug"
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected %q to contain %q", err.Error(), expected)
	}
}

func TestParseSlug_goodString(t *testing.T) {
	user, name, err := parseSlug("hashicorp/project")
	if err != nil {
		t.Fatal(err)
	}

	if user != "hashicorp" {
		t.Fatalf("expected %q to be %q", user, "hashicorp")
	}

	if name != "project" {
		t.Fatalf("expected %q to be %q", name, "project")
	}
}

func TestHarmonyClient_noURL(t *testing.T) {
	client, err := harmonyClient(&UploadOpts{})
	if err != nil {
		t.Fatal(err)
	}

	expected := harmony.DefaultClient()
	if !reflect.DeepEqual(client, expected) {
		t.Fatalf("expected %q to be %q", client, expected)
	}
}

func TestHarmonyClient_customURL(t *testing.T) {
	url := "https://harmony.company.com"
	client, err := harmonyClient(&UploadOpts{URL: url})
	if err != nil {
		t.Fatal(err)
	}

	expected, err := harmony.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(client, expected) {
		t.Fatalf("expected %q to be %q", client, expected)
	}
}

func TestHarmonyClient_token(t *testing.T) {
	token := "abcd1234"
	client, err := harmonyClient(&UploadOpts{Token: token})
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
