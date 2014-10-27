package main

import (
	"io/ioutil"
	"testing"
)

func tempFile(t *testing.T) string {
	tf, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer tf.Close()

	return tf.Name()
}
