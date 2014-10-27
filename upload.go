package main

import (
	"io"
)

// UploadOpts are the options for uploading the archive.
type UploadOpts struct {
	App   string
	Token string
}

// Upload uploads the reader, representing a single archive, to the
// application given by UploadOpts.
//
// The Upload happens async and the return values are the done channel,
// the error channel, and then an error that can happen during initialization.
// If error is returned, then the channels will be nil and the upload never
// started. Otherwise, the upload has started in the background and is not
// done until the done channel or error channel send a value. Once either send
// a value, the upload is stopped.
func Upload(r io.Reader, opts *UploadOpts) (<-chan struct{}, <-chan error, error) {
	return nil, nil, nil
}
