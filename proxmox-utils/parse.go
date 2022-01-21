package pu

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var formatRx = regexp.MustCompile(`format=(cloop|cow|qcow|qcow2|qed|raw|vmdk)`)
var storageKeyRx = regexp.MustCompile(`^(ide|sata|scsi)([0-9]+)$`)

func parseVMInfo(cfg string) (ret VM, err error) {
	for _, row := range strings.Split(cfg, "\n") {
		if !strings.ContainsRune(row, ':') {
			continue
		}

		parts := strings.SplitN(row, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "name":
			ret.Name = value
		case "scsihw":
			ret.SCSIHardware = value
		case "cores":
			ret.Cores, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ret, fmt.Errorf("error decoding cores number: %w", err)
			}
		case "memory":
			ret.Memory, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ret, fmt.Errorf("error decoding memory number: %w", err)
			}
		case "sockets":
			ret.Sockets, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ret, fmt.Errorf("error decoding sockets number: %w", err)
			}
		default:
			switch {
			case strings.HasPrefix(key, "net"):
				netID, err := strconv.ParseInt(strings.ReplaceAll(key, "net", ""), 10, 32)
				if err != nil {
					return ret, fmt.Errorf("error decoding net number: %w", err)
				}
				components := strings.Split(value, "=")
				ret.Networks = append(ret.Networks, Network{ID: netID, Card: components[0]})
			case strings.HasPrefix(key, DiskBusIDE):
				fallthrough
			case strings.HasPrefix(key, DiskBusSATA):
				fallthrough
			case strings.HasPrefix(key, DiskBusSCSI):
				bus, slot, err := parseStorageKey(key)
				if err != nil {
					return ret, fmt.Errorf("error decoding storage key: %w", err)
				}
				volume, isCdrom, format := parseStorageLine(value)
				if volume == "" {
					continue
				}

				ret.Disks = append(ret.Disks, Disk{
					ID:     slot,
					Bus:    bus,
					Volume: volume,
					Format: format,
					Cdrom:  isCdrom,
				})
			}
		}
	}
	return ret, nil
}

func parseStorageKey(key string) (bus string, slot int64, err error) {
	if !storageKeyRx.MatchString(key) {
		return "", 0, fmt.Errorf("storage key error")
	}

	split := storageKeyRx.FindAllStringSubmatch(key, -1)
	if len(split) != 1 || len(split[0]) != 3 {
		return "", 0, fmt.Errorf("storage key error")
	}
	bus = split[0][1]
	slot, err = strconv.ParseInt(split[0][2], 10, 32)
	return
}

func parseStorageLine(value string) (volume string, isCdrom bool, format string) {
	isCdrom = strings.Contains(value, "media=cdrom")
	components := strings.Split(value, ",")

	if components[0] != "none" {
		volume = components[0]
	}

	if formatRx.MatchString(value) {
		formatvalues := formatRx.FindAllStringSubmatch(value, -1)
		if len(formatvalues) == 1 && len(formatvalues[0]) == 2 {
			format = formatvalues[0][1]
		}
	}
	return
}
