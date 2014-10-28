package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun__versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("harmony-upload -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %s to eq %s", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("harmony-upload v%s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestRun_parseError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("harmony-upload -bacon delicious", " ")

	status := cli.Run(args)
	if status != ExitCodeParseFlagsError {
		t.Errorf("expected %s to eq %s", status, ExitCodeParseFlagsError)
	}

	expected := "flag provided but not defined: -bacon"
	if !strings.Contains(errStream.String(), expected) {
		t.Fatalf("expected %q to contain %q", errStream.String(), expected)
	}
}

func TestRun_includeFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("harmony-upload -include foo hashicorp/project .", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %s to eq %s", status, ExitCodeOK)
	}
}

func TestRun_excludeFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("harmony-upload -exclude bar hashicorp/project .", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %s to eq %s", status, ExitCodeOK)
	}
}
