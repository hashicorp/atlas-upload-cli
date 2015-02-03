package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/atlas-go/archive"
	"github.com/hashicorp/logutils"
	"github.com/mitchellh/ioprogress"
)

// Exit codes are int values that represent an exit code for a particular error.
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

// levelFilter is the log filter with pre-defined levels
var levelFilter = &logutils.LevelFilter{
	Levels: []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERR"},
}

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the standard out and standard error streams to
	// write messages from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments. The first argument is always
// the name of the application. This method slices accordingly.
func (cli *CLI) Run(args []string) int {
	// Initialize the logger to start (overridden later if debug is given)
	cli.initLogger(os.Getenv("ATLAS_LOG"))

	var debug, version bool
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
	flags.Var((*FlagSliceVar)(&archiveOpts.Exclude), "exclude",
		"files/folders to exclude")
	flags.Var((*FlagSliceVar)(&archiveOpts.Include), "include",
		"files/folders to include")
	flags.Var((*FlagMetadataVar)(&uploadOpts.Metadata), "metadata",
		"arbitrary metadata to pass along with the request")
	flags.BoolVar(&debug, "debug", false,
		"turn on debug output")
	flags.BoolVar(&version, "version", false,
		"display the version")

	// Parse all the flags
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	// Turn on debug mode if requested
	if debug {
		levelFilter.SetMinLevel(logutils.LogLevel("DEBUG"))
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
		fmt.Fprintf(cli.errStream, "cli: must specify two arguments - app, path\n")
		flags.Usage()
		return ExitCodeBadArgs
	}

	// Get the name of the app and the path to archive
	slug, path := parsedArgs[0], parsedArgs[1]
	uploadOpts.Slug = slug

	// Get the archive reader
	r, err := archive.CreateArchive(path, &archiveOpts)
	if err != nil {
		fmt.Fprintf(cli.errStream, "error archiving: %s\n", err)
		return ExitCodeArchiveError
	}
	defer r.Close()

	// Put a progress bar around the reader
	pr := &ioprogress.Reader{
		Reader: r,
		Size:   r.Size,
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, func(p, t int64) string {
			return fmt.Sprintf(
				"Uploading %s: %s",
				slug,
				ioprogress.DrawTextFormatBytes(p, t))
		}),
	}

	// Start the upload
	doneCh, uploadErrCh, err := Upload(pr, r.Size, &uploadOpts)
	if err != nil {
		fmt.Fprintf(cli.errStream, "error starting upload: %s\n", err)
		return ExitCodeUploadError
	}

	select {
	case err := <-uploadErrCh:
		fmt.Fprintf(cli.errStream, "error uploading: %s\n", err)
		return ExitCodeUploadError
	case version := <-doneCh:
		fmt.Printf("Uploaded %s v%d\n", slug, version)
	}

	return ExitCodeOK
}

// initLogger gets the log level from the environment, falling back to DEBUG if
// nothing was given.
func (cli *CLI) initLogger(level string) {
	minLevel := strings.ToUpper(strings.TrimSpace(level))
	if minLevel == "" {
		minLevel = "WARN"
	}

	levelFilter.Writer = cli.errStream
	levelFilter.SetMinLevel(logutils.LogLevel(level))
	log.SetOutput(levelFilter)
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

  -exclude=<path>     Glob pattern of files or directories to exlude (this may
                      be specified multiple times)
  -include=<path>     Glob pattern of files/directories to include (this may be
                      specified multiple times, any excludes will override
                      conflicting includes)
  -address=<url>      The address of the Atlas server
  -token=<token>      The Atlas API token
  -vcs                Use VCS to determine which files to include/exclude

  -metadata<k=v>      Arbitrary key-value (string) metadata to be sent with the
                      upload; may be specified multiple times

  -debug              Turn on debug output
  -version            Print the version of this application
`
