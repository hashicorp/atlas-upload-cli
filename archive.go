package main

import (
	"io"
)

// ArchiveOpts are the options for defining how the archive will be built.
type ArchiveOpts struct {
	Exclude []string
	Include []string
	VCS     bool
}

// IsSet says whether any options were set.
func (o *ArchiveOpts) IsSet() bool {
	return len(o.Exclude) > 0 || len(o.Include) > 0 || o.VCS
}

// Archive takes the given path and ArchiveOpts and archives it.
//
// The archive is done async and streamed via the io.ReadCloser returned.
// The reader is blocking: data is only compressed and written as data is
// being read from the reader. Because of this, any user doesn't have to
// worry about quickly reading data to avoid memory bloat.
//
// The archive can be read with the io.ReadCloser that is returned. The error
// returned is an error that happened before archiving started, so the
// ReadCloser doesn't need to be closed (and should be nil). The error
// channel are errors that can happen while archiving is happening. When
// an error occurs on the channel, reading should stop and be closed.
func Archive(
	path string, opts *ArchiveOpts) (io.ReadCloser, <-chan error, error) {
	return nil, nil, nil
}
