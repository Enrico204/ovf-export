package main

import (
	"fmt"
	pu "gitlab.com/enrico204/ovf-export/proxmox-utils"
	qu "gitlab.com/enrico204/ovf-export/qemu-utils"
	"os"
	"os/exec"
	"path"
	"strings"
)

type conversionResult struct {
	FileName             string
	TotalCapacity        int64
	SparseFileOccupation int64
	FileSize             int64
}

func importDisk(file string, outdir string, bus string) (conversionResult, error) {
	var ret conversionResult

	if strings.HasSuffix(file, ".iso") {
		ret.FileName = outdir + path.Base(file)

		// I'm quite lazy today, so we copy using cp
		_, err := exec.Command("cp", file, outdir).Output()
		if err != nil {
			return ret, fmt.Errorf("error copying file: %w", err)
		}
	} else {
		ret.FileName = outdir + path.Base(strings.ReplaceAll(file, ".qcow2", ".vmdk"))

		var adapterType = qu.AdapterTypeIDE
		if bus == pu.DiskBusSCSI || bus == pu.DiskBusSATA {
			adapterType = qu.AdapterTypeLSILogic
		}

		err := qu.ConvertDisk(file, ret.FileName, adapterType, false, qu.SubFormatStreamOptimized, "", nil, false)
		if err != nil {
			return ret, fmt.Errorf("error converting disk: %w", err)
		}
	}

	diskInfo, err := qu.GetDiskInfo(ret.FileName)
	if err != nil {
		return ret, err
	}

	fstat, err := os.Stat(ret.FileName)
	if err != nil {
		return ret, err
	}

	ret.FileSize = fstat.Size()
	ret.TotalCapacity = diskInfo.VirtualSize
	ret.SparseFileOccupation = diskInfo.DiskSize

	return ret, nil
}
