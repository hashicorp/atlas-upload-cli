package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	fi, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}

	// Direct file paths cannot have archive options
	if !fi.IsDir() && opts.IsSet() {
		return nil, nil, fmt.Errorf(
			"Options such as exclude, include, and VCS can't be set when " +
				"the path is a file.")
	}

	if fi.IsDir() {
		return archiveDir(path, opts)
	} else {
		return archiveFile(path, opts)
	}
}

func archiveFile(
	path string, opts *ArchiveOpts) (io.ReadCloser, <-chan error, error) {
	// TODO: if file is already gzipped, then send it along
	// TODO: if file is not gzipped, then... error? or do we tar + gzip?

	return nil, nil, nil
}

func archiveDir(
	root string, opts *ArchiveOpts) (io.ReadCloser, <-chan error, error) {
	var vcsInclude []string
	if opts.VCS {
		var err error
		vcsInclude, err = vcsFiles(root)
		if err != nil {
			return nil, nil, err
		}
	}

	// We're going to write to an io.Pipe so that we can ensure the other
	// side is reading as we're writing.
	pr, pw := io.Pipe()

	// Buffer the writer so that we can keep some data moving in memory
	// while we're compressing. 4M should be good.
	bufW := bufio.NewWriterSize(pw, 4096 * 1024)

	// Gzip compress all the output data
	gzipW := gzip.NewWriter(bufW)

	// Build the function that'll do all the compression
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path from the path since it contains the root
		// plus the path.
		subpath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// If we have a list of VCS files, check that first
		skip := false
		if len(vcsInclude) > 0 {
			skip = true
			for _, f := range vcsInclude {
				if f == subpath {
					skip = false
					break
				}
			}
		}

		// TODO: include/exclude lists

		// If we have to skip this file, then skip it, properly skipping
		// children if we're a directory.
		if skip {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		// Read the symlink target. We don't track the error because
		// it doesn't matter if there is an error.
		target, _ := os.Readlink(path)

		// Build the file header for the tar entry
		header, err := tar.FileInfoHeader(info, target)
		if err != nil {
			return fmt.Errorf(
				"Failed creating archive header: %s", path)
		}

		tarW := tar.NewWriter(gzipW)
		defer tarW.Close()

		// Write the header first to the archive.
		if err := tarW.WriteHeader(header); err != nil {
			return fmt.Errorf(
				"Failed writing archive header: %s", path)
		}

		// If it is a directory, then we're done (no body to write)
		if info.IsDir() {
			return nil
		}

		// Open the target file to write the data
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf(
				"Failed opening file '%s' to write compressed archive.", path)
		}
		defer f.Close()

		if _, err = io.Copy(tarW, f); err != nil {
			return fmt.Errorf(
				"Failed copying file to archive: %s", path)
		}

		return nil
	}

	// Create all our channels so we can send data through some tubes
	// to other goroutines.
	errCh := make(chan error, 1)
	go func() {
		err := filepath.Walk(root, walkFn)

		// TODO: errors for everything

		// Close the gzip writer
		gzipW.Close()

		// Flush the buffer
		bufW.Flush()

		// Close the pipe
		pw.Close()

		// Send any error we might have down the pipe if we have one
		if err != nil {
			errCh <- err
		}
	}()

	return pr, errCh, nil
}
