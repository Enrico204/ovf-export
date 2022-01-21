package lu

import (
	"encoding/hex"
	"fmt"
	"github.com/digitalocean/go-libvirt"
)

func DecodeLibvirtUUID(uuid string) (libvirt.UUID, error) {
	var ret libvirt.UUID
	id, err := hex.DecodeString(uuid)
	if err != nil {
		return ret, fmt.Errorf("failed to decode domain UUID: %w", err)
	} else if len(id) != 16 {
		return ret, fmt.Errorf("failed to decode domain UUID")
	}

	copy(ret[:], id)
	return ret, nil
}
