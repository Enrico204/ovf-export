package ova

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func GenerateFromOVFDir(ovfdir string) error {
	ovaFilePath := path.Join(ovfdir, "..", path.Base(ovfdir)+".ova")
	_ = os.Remove(ovaFilePath)

	entries, err := os.ReadDir(ovfdir)
	if err != nil {
		return fmt.Errorf("can't read OVF dir: %w", err)
	}

	var cmd = []string{"cf", ovaFilePath, "--format=ustar", "-C", ovfdir}
	var ovfFile string
	var manifestFile string
	var otherFiles []string

	// Generate OVA: first the OVF file, then others
	for _, entry := range entries {
		switch {
		case entry.IsDir():
			continue
		case strings.HasSuffix(entry.Name(), ".ovf"):
			ovfFile = entry.Name()
		case strings.HasSuffix(entry.Name(), ".mf"):
			manifestFile = entry.Name()
		default:
			otherFiles = append(otherFiles, entry.Name())
		}
	}

	if ovfFile == "" {
		return fmt.Errorf("missing OVF file")
	}
	cmd = append(cmd, ovfFile)

	if manifestFile == "" {
		cmd = append(cmd, manifestFile)
	}

	_, err = exec.Command("tar", append(cmd, otherFiles...)...).CombinedOutput() //nolint:gosec
	if err != nil {
		return fmt.Errorf("can't create OVA: %w", err)
	}
	return nil
}
