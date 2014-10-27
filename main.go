package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var archiveOpts ArchiveOpts
	var uploadOpts UploadOpts

	flag.BoolVar(&archiveOpts.VCS, "vcs", false, "vcs")
	flag.StringVar(&uploadOpts.Token, "token", "", "token")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Must specify one argument: app path\n")
		usage()
		return 1
	}

	// Get the archive reader.
	r, archiveErrCh, err := Archive(args[1], &archiveOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with path: %s\n", err)
		return 1
	}
	defer r.Close()

	// Start the upload
	doneCh, uploadErrCh, err := Upload(r, &uploadOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting upload: %s\n", err)
		return 1
	}

	select {
	case err := <-archiveErrCh:
		fmt.Fprintf(os.Stderr, "Error archiving: %s\n", err)
		return 1
	case err := <-uploadErrCh:
		fmt.Fprintf(os.Stderr, "Error uploading: %s\n", err)
		return 1
	case <-doneCh:
	}

	return 0
}

func usage() {
	fmt.Fprintf(os.Stderr, usageStr, os.Args[0])
}

const usageStr = `
Usage: %s [options] app path

  Upload application code or artifacts to Harmony for initiating deployments.

  "app" is the name of the application to upload to within Harmony.

  If path is a directory, it will be compressed (gzip tar) and uploaded
  in its entirety. The root of the archive will be the path. For clarify:
  if you upload the "foo/" directory, then the file "foo/version" will be
  "version" in the archive since "foo/" is the root.

  A path must be specified. Due to the nature of this application, it does
  not default to using the current working directory automatically.

Options:

  -exclude=<path>          Glob pattern of files or directories to exlude.
                           This can be specified multiple times.
  -include=<path>          Glob pattern of files to directories to include.
                           This can be specified multiple times. Any excludes
                           will override any conflicting includes.
  -token=<token>           Harmony API token.
  -vcs                     If path is version controlled, it will use the VCS
                           to determine what files to include/exclude.
`
