//nolint:forbidigo
package main

import (
	"fmt"
	"gitlab.com/enrico204/ovf-export/ovf/manifest"
	"os"
	"regexp"
)

var replaceLastOvfRx = regexp.MustCompile(`\.ovf$`)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <path-to-ovf-file>", os.Args[0])
		os.Exit(1)
	}

	var ovfFile = os.Args[1]
	var mfFile = replaceLastOvfRx.ReplaceAllString(ovfFile, ".mf")

	// Parse/dump manifest
	manifestContent, err := manifest.ParseFile(mfFile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	} else if err == nil {
		fmt.Println(manifestContent)
	}
}
