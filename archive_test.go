package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestArchiveOptsIsSet(t *testing.T) {
	cases := []struct {
		Opts *ArchiveOpts
		Set  bool
	}{
		{
			&ArchiveOpts{},
			false,
		},
		{
			&ArchiveOpts{VCS: true},
			true,
		},
		{
			&ArchiveOpts{Exclude: make([]string, 0, 0)},
			false,
		},
		{
			&ArchiveOpts{Exclude: []string{"foo"}},
			true,
		},
		{
			&ArchiveOpts{Include: make([]string, 0, 0)},
			false,
		},
		{
			&ArchiveOpts{Include: []string{"foo"}},
			true,
		},
	}

	for i, tc := range cases {
		if tc.Opts.IsSet() != tc.Set {
			t.Fatalf("%d: expected %#v", i, tc.Set)
		}
	}
}

func TestArchive_fileNoExist(t *testing.T) {
	tf := tempFile(t)
	if err := os.Remove(tf); err != nil {
		t.Fatalf("err: %s", err)
	}

	r, errCh, err := Archive(tf, nil)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	if r != nil {
		t.Fatal("should be nil")
	}
	if errCh != nil {
		t.Fatal("should be nil")
	}
}

func TestArchive_fileWithOpts(t *testing.T) {
	r, errCh, err := Archive(tempFile(t), &ArchiveOpts{VCS: true})
	if err == nil {
		t.Fatal("err should not be nil")
	}
	if r != nil {
		t.Fatal("should be nil")
	}
	if errCh != nil {
		t.Fatal("should be nil")
	}
}

func TestArchive_dirNoVCS(t *testing.T) {
	r, errCh, err := Archive(testFixture("archive-flat"), new(ArchiveOpts))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{
		"baz.txt",
		"foo.txt",
	}

	entries := testArchive(t, r, errCh)
	if !reflect.DeepEqual(entries, expected) {
		t.Fatalf("bad: %#v", entries)
	}
}

func testArchive(t *testing.T, r io.ReadCloser, errCh <-chan error) []string {
	// Finish the archiving process in-memory
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("err: %s", err)
	}
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	default:
	}

	gzipR, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tarR := tar.NewReader(gzipR)

	// Read all the entries
	result := make([]string, 0, 5)
	for {
		hdr, err := tarR.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		result = append(result, hdr.Name)
	}

	sort.Strings(result)
	return result
}
