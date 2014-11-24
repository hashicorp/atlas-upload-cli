package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/atlas-go/archive"
)

// Exit codes are int valuse that represent an exit code for a particular error.
// Sub-systems may check this unique error to determine the cause of an error
// without parsing the output or help text.
const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError
	ExitCodeBadArgs
	ExitCodeArchiveError
	ExitCodeUploadError
)

// CLI is the command line object
type CLI struct {
	// outSteam and errStream are the standard out and standard error streams to
	// write messages from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments. The first arugment is always
// the name of the application. This method slices accordingly.
func (cli *CLI) Run(args []string) int {
	var version bool
	var archiveOpts archive.ArchiveOpts
	var uploadOpts UploadOpts

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprintf(cli.errStream, usage, Name)
	}
	flags.BoolVar(&archiveOpts.VCS, "vcs", false,
		"use VCS to detect which files to upload")
	flags.StringVar(&uploadOpts.URL, "address", "",
		"Atlas server address")
	flags.StringVar(&uploadOpts.Token, "token", "",
		"Atlas API token")
	flags.Var((*flagSliceVar)(&archiveOpts.Exclude), "exclude",
		"files/folders to exclude")
	flags.Var((*flagSliceVar)(&archiveOpts.Include), "include",
		"files/folders to include")
	flags.BoolVar(&version, "version", false,
		"display the version")

	// Parse all the flags
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	// Version
	if version {
		fmt.Fprintf(cli.errStream, "%s v%s\n", Name, Version)
		return ExitCodeOK
	}

	// Get the parsed arguments (the ones left over after all the flags have been
	// parsed)
	parsedArgs := flags.Args()

	if len(parsedArgs) != 2 {
		fmt.Fprintf(cli.errStream, "cli: must specify two arguments - app, path")
		flags.Usage()
		return ExitCodeBadArgs
	}

	// Get the name of the app and the path to archive
	slug, path := parsedArgs[0], parsedArgs[1]
	uploadOpts.Slug = slug

	// Get the archive reader
	r, err := archive.CreateArchive(path, &archiveOpts)
	if err != nil {
		fmt.Fprintf(cli.errStream, "error archiving: %s", err)
		return ExitCodeArchiveError
	}
	defer r.Close()

	// Start the upload
	doneCh, uploadErrCh, err := Upload(r, r.Size, &uploadOpts)
	if err != nil {
		fmt.Fprintf(cli.errStream, "error starting upload: %s", err)
		return ExitCodeUploadError
	}

	select {
	case err := <-uploadErrCh:
		fmt.Fprintf(cli.errStream, "error uploading: %s", err)
		return ExitCodeUploadError
	case <-doneCh:
	}

	return ExitCodeOK
}

const usage = `
Usage: %s [options] app path

  Upload application code or artifacts to Atlas for initiating deployments.

  "app" is the name of the application to upload to within Atlas.

  If path is a directory, it will be compressed (gzip tar) and uploaded
  in its entirety. The root of the archive will be the path. For clarity:
  if you upload the "foo/" directory, then the file "foo/version" will be
  "version" in the archive since "foo/" is the root.

  A path must be specified. Due to the nature of this application, it does
  not default to using the current working directory automatically.

Options:

  -exclude=<path>     Glob pattern of files or directories to exlude. This can
                      be specified multiple times.
  -include=<path>     Glob pattern of files/directories to include. This can be
                      specified multiple times. Any excludes will override
                      conflicting includes.
  -address=<url>      The address of the Atlas server
  -token=<token>      The Atlas API token
  -vcs                Use VCS to determine which files to include/exclude

  -version            Print the version of this application
`

// flagSliceVar is a special flag that permits the value to be supplied more
// than once. Values are pushed onto a string slice.
type flagSliceVar []string

func (fsv *flagSliceVar) String() string {
	return strings.Join(*fsv, ",")
}

func (fsv *flagSliceVar) Set(value string) error {
	if *fsv == nil {
		*fsv = make([]string, 0, 1)
	}
	*fsv = append(*fsv, value)
	return nil
}
