package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

const fixturesDir = "./test-fixtures"

func tempFile(t *testing.T) string {
	tf, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer tf.Close()

	return tf.Name()
}

func testFixture(n string) string {
	return filepath.Join(fixturesDir, n)
}
