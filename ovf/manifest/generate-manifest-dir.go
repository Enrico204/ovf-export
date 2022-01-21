package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func GenerateManifestFromOVFDir(outdir string) error {
	entries, err := os.ReadDir(outdir)
	if err != nil {
		return err
	}

	var ovfname = ""

	// Generate manifest
	var mf = make(Content, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".mf") {
			continue
		}

		err = mf.AddFile(path.Join(outdir, entry.Name()))
		if err != nil {
			return err
		}

		if strings.HasSuffix(entry.Name(), ".ovf") {
			ovfname = entry.Name()
		}
	}

	if ovfname == "" {
		return fmt.Errorf("missing OVF file")
	}

	manifestContent, err := mf.Build()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(outdir, strings.ReplaceAll(ovfname, ".ovf", ".mf")), manifestContent, 0600)
}
