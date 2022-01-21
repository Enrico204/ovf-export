package qu

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var diskSizeRx = regexp.MustCompile(`^\s*([0-9.]+)\s*([KMGT]i?B)\s*(\(([0-9]+) bytes\))?\s*$`)

type DiskInfo struct {
	FileFormat  string
	VirtualSize int64
	DiskSize    int64
}

func GetDiskInfo(filePath string) (DiskInfo, error) {
	var ret DiskInfo

	out, err := exec.Command("qemu-img", "info", filePath).Output()
	if err != nil {
		return ret, err
	}

	for _, row := range strings.Split(string(out), "\n") {
		if !strings.ContainsRune(row, ':') {
			continue
		}

		parts := strings.Split(row, ":")
		if len(parts) != 2 {
			return ret, fmt.Errorf("invalid output from qemu-img")
		}

		switch {
		case strings.HasPrefix(row, "file format"):
			ret.FileFormat = strings.TrimSpace(parts[1])
		case strings.HasPrefix(row, "virtual size"):
			ret.VirtualSize, err = sizeToBytes(parts[1])
			if err != nil {
				return ret, err
			}
		case strings.HasPrefix(row, "disk size"):
			ret.DiskSize, err = sizeToBytes(parts[1])
			if err != nil {
				return ret, err
			}
		}
	}
	return ret, nil
}

func sizeToBytes(txt string) (int64, error) {
	if !diskSizeRx.MatchString(txt) {
		return 0, fmt.Errorf("invalid qemu-img info size format, no match")
	}
	values := diskSizeRx.FindAllStringSubmatch(txt, -1)
	if len(values) != 1 || len(values[0]) != 5 {
		return 0, fmt.Errorf("invalid qemu-img info size format, wrong match")
	}

	if values[0][4] != "" {
		// The size in bytes is already available, skip conversion
		return strconv.ParseInt(values[0][4], 10, 64)
	}

	humanSize, err := strconv.ParseFloat(values[0][1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid qemu-img info size format: %w", err)
	}

	switch values[0][2] {
	case "KiB":
		return int64(humanSize * 1024), nil
	case "MiB":
		return int64(humanSize * 1024 * 1024), nil
	case "GiB":
		return int64(humanSize * 1024 * 1024 * 1024), nil
	case "TiB":
		return int64(humanSize * 1024 * 1024 * 1024 * 1024), nil
	case "B":
		return int64(humanSize), nil
	default:
		return 0, fmt.Errorf("invalid qemu-img info size unit")
	}
}
