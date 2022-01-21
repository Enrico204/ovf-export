package pu

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Disk struct {
	ID     int64
	Bus    string
	Volume string
	Format string
	Cdrom  bool
}

func (d *Disk) GetFilePath() (string, error) {
	out, err := exec.Command("pvesm", "path", d.Volume).CombinedOutput() //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("error getting file path: %w %s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

type Network struct {
	ID   int64
	Card string
}

type VM struct {
	Cores        int64
	Memory       int64
	Sockets      int64
	Name         string
	SCSIHardware string
	Disks        []Disk
	Networks     []Network
}

func ExportVMInfo(vmID int64) (VM, error) {
	out, err := exec.Command("qm", "config", strconv.FormatInt(vmID, 10)).CombinedOutput() //nolint:gosec
	if err != nil {
		return VM{}, fmt.Errorf("error getting VM info: %w %s", err, out)
	}
	return parseVMInfo(string(out))
}
