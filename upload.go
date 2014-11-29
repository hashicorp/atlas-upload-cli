package main

import (
	"fmt"
	"io"

	"github.com/hashicorp/atlas-go/v1"
)

// UploadOpts are the options for uploading the archive.
type UploadOpts struct {
	// URL is the Atlas endpoint. If this value is not specified, the uploader
	// will default to the public Atlas install as defined in the atlas-go
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
func Upload(r io.Reader, size int64, opts *UploadOpts) (<-chan struct{}, <-chan error, error) {
	// Create the client
	client, err := atlasClient(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("upload: %s", err)
	}

	// Separate the slug into the user and name components
	user, name, err := atlas.ParseSlug(opts.Slug)
	if err != nil {
		return nil, nil, fmt.Errorf("upload: %s", err)
	}

	// Get the app
	app, err := client.App(user, name)
	if err != nil {
		if err == atlas.ErrNotFound {
			// Application doesn't exist, attempt to create it
			app, err = client.CreateApp(user, name)
		}

		if err != nil {
			return nil, nil, fmt.Errorf("upload: %s", err)
		}
	}

	doneCh, errCh := make(chan struct{}), make(chan error)

	// Start the upload
	go process(func() error {
		return client.UploadApp(app, r, size)
	}, doneCh, errCh)

	return doneCh, errCh, nil
}

// Create the client - if a URL is given, construct a new Client from the URL,
// but if not URL is given, use the default client.
func atlasClient(opts *UploadOpts) (*atlas.Client, error) {
	var client *atlas.Client
	var err error

	if opts.URL == "" {
		client = atlas.DefaultClient()
	} else {
		client, err = atlas.NewClient(opts.URL)
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
