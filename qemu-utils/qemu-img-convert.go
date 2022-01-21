package qu

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

//nolint:deadcode
const (
	AdapterTypeIDE       = "ide"
	AdapterTypeLSILogic  = "lsilogic"
	AdapterTypeBusLogic  = "buslogic"
	AdapterTypeLegacyESX = "legacyESX"

	SubFormatMonolithicSparse     = "monolithicSparse"
	SubFormatMonolithicFlat       = "monolithicFlat"
	SubFormatTwoGbMaxExtentSparse = "twoGbMaxExtentSparse"
	SubFormatTwoGbMaxExtentFlat   = "twoGbMaxExtentFlat"
	SubFormatStreamOptimized      = "streamOptimized"
)

func ConvertDisk(src string, dst string, adapterType string, vmdkv6 bool, subformat string, backingFile string, size *uint64, zeroedGrain bool) error {
	var vmdkopts = make([]string, 0)

	if adapterType != "" {
		vmdkopts = append(vmdkopts, "adapter_type="+adapterType)
	}
	if subformat != "" {
		vmdkopts = append(vmdkopts, "subformat="+subformat)
	}
	if backingFile != "" {
		vmdkopts = append(vmdkopts, "backing_file="+backingFile)
	}
	if size != nil {
		vmdkopts = append(vmdkopts, "size="+strconv.Itoa(int(*size)))
	}
	if vmdkv6 {
		vmdkopts = append(vmdkopts, "compat6")
	}
	if zeroedGrain {
		vmdkopts = append(vmdkopts, "zeroed_grain")
	}

	var args = []string{"convert", "-m", strconv.Itoa(runtime.NumCPU()), "-O", "vmdk"}

	if len(vmdkopts) > 0 {
		args = append(args, "-o", strings.Join(vmdkopts, ","))
	}

	out, err := exec.Command("qemu-img", append(args, src, dst)...).CombinedOutput() //nolint:gosec
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}
	return nil
}
