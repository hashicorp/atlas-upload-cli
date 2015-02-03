package main

import (
	"fmt"
	"strings"
)

// FlagMetadataVar is a flag.Value implementation for parsing user variables
// from the command-line in the format of 'key=value'.
type FlagMetadataVar map[string]interface{}

func (v *FlagMetadataVar) String() string {
	return ""
}

func (v *FlagMetadataVar) Set(raw string) error {
	idx := strings.Index(raw, "=")
	if idx == -1 {
		return fmt.Errorf("Missing '=' in argument: %s", raw)
	}

	if *v == nil {
		*v = make(map[string]interface{})
	}

	key, value := raw[0:idx], raw[idx+1:]
	(*v)[key] = value

	return nil
}

// FlagSliceVar is a special flag that permits the value to be supplied more
// than once. Values are pushed onto a string slice.
type FlagSliceVar []string

func (fsv *FlagSliceVar) String() string {
	return strings.Join(*fsv, ",")
}

func (fsv *FlagSliceVar) Set(value string) error {
	if *fsv == nil {
		*fsv = make([]string, 0, 1)
	}
	*fsv = append(*fsv, value)
	return nil
}
