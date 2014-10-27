package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// VCS is a struct that explains how to get the file list for a given
// VCS.
type VCS struct {
	Name string

	// Detect is a list of files/folders that if they exist, signal that
	// this VCS is the VCS in use.
	Detect []string

	// Files returns the files that are under version control for the
	// given path.
	Files func(path string) ([]string, error)
}

// VCSList is the list of VCS we recognize.
var VCSList = []*VCS{
	&VCS{
		Name:   "git",
		Detect: []string{".git/"},
		Files:  vcsFilesCmd("git", "ls-files"),
	},
}

// vcsFilesCmd creates a Files-compatible function that reads the files
// by executing the command in the repository path and returning each
// line in stdout.
func vcsFilesCmd(args ...string) func(string) ([]string, error) {
	return func(path string) ([]string, error) {
		var stderr, stdout bytes.Buffer

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = path
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf(
				"Error executing %s: %s",
				strings.Join(args, " "),
				err)
		}

		// Read each line of output as a path
		result := make([]string, 0, 100)
		scanner := bufio.NewScanner(&stdout)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}

		return result, nil
	}
}
