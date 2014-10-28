package main

import (
	"fmt"
	"io"

	harmony "github.com/hashicorp/harmony-go"
)

// UploadOpts are the options for uploading the archive.
type UploadOpts struct {
	// URL is the Harmony endpoint. If this value is not specified, the uploader
	// will default to the public Harmony install as defined in the harmony-go
	// client.
	URL string

	// Slug is the "user/name" of the application to upload.
	Slug string

	// Token is the API token to upload with.
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
	// Create the client
	client, err := harmonyClient(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("upload: %s", err)
	}

	// Separate the slug into the user and name components
	user, name, err := harmony.ParseSlug(opts.Slug)
	if err != nil {
		return nil, nil, fmt.Errorf("upload: %s", err)
	}

	// Get the app
	app, err := client.App(user, name)
	if err != nil {
		return nil, nil, fmt.Errorf("upload: %s", err)
	}

	doneCh, errCh := make(chan struct{}), make(chan error)

	// Start the upload
	go process(func() error {
		return client.UploadApp(app, r)
	}, doneCh, errCh)

	return doneCh, errCh, nil
}

// Create the client - if a URL is given, construct a new Client from the URL,
// but if not URL is given, use the default client.
func harmonyClient(opts *UploadOpts) (*harmony.Client, error) {
	var client *harmony.Client
	var err error

	if opts.URL == "" {
		client = harmony.DefaultClient()
	} else {
		client, err = harmony.NewClient(opts.URL)
	}

	if opts.Token != "" {
		client.Token = opts.Token
	}

	return client, err
}

// process takes an arbitrary function that returns an error, a doneCh, and an
// errCh. The function is executed in serial and any errors are pushed onto the
// errCh. This function blocks until it finishes, so it should be run from a
// separate goroutine.
func process(f func() error, doneCh chan<- struct{}, errCh chan<- error) {
	if err := f(); err != nil {
		errCh <- fmt.Errorf("upload: %s", err)
		return
	}
	close(doneCh)
}
