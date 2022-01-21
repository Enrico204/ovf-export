//nolint:forbidigo
package main

import (
	"fmt"
	v1 "gitlab.com/enrico204/ovf-export/ovf/v1"
)

func main() {
	var zero = 0
	var tre = 3

	var ovf = v1.Envelope{
		XMLLang: "en-US",
		References: v1.References{
			Files: []v1.File{
				{
					ID:   "file1",
					Href: "disk1.vmdk",
				},
			},
		},
		DiskSection: v1.DiskSection{
			Disks: []v1.Disk{{
				Capacity: 1234,
				DiskID:   "vmdisk1",
				FileRef:  "file1",
				Format:   v1.DiskFormatVMDKStreamOptimized,
			}},
		},
		NetworkSection: v1.NetworkSection{
			Networks: []v1.Network{{
				Name: "NAT",
			}},
		},
		VirtualSystem: v1.VirtualSystem{
			ID: "Mikrotik",
			OperatingSystemSection: v1.OperatingSystemSection{
				ID:          "0",
				Description: "Other",
			},
			VirtualHardwareSection: v1.VirtualHardwareSection{
				System: v1.System{
					ElementName:             "Virtual Hardware Family",
					InstanceID:              0,
					VirtualSystemIdentifier: "Mikrotik",
					VirtualSystemType:       "virtualbox-2.2",
				},
				Items: []v1.Item{
					{
						Caption:         "1 virtual CPU",
						Description:     "Number of virtual CPUs",
						ElementName:     "1 virtual CPU",
						InstanceID:      1,
						ResourceType:    3,
						VirtualQuantity: 1,
					},
					{
						AllocationUnits: "MegaBytes",
						Caption:         "1024 MB of memory",
						Description:     "Memory Size",
						ElementName:     "1024 MB of memory",
						InstanceID:      2,
						ResourceType:    4,
						VirtualQuantity: 1024,
					},
					{
						Address:         &zero,
						Caption:         "ideController0",
						Description:     "IDE Controller",
						ElementName:     "ideController0",
						InstanceID:      3,
						ResourceType:    5,
						ResourceSubType: "PIIX4",
					},
					{
						AddressOnParent: &zero,
						Caption:         "disk1",
						Description:     "Disk Image",
						ElementName:     "disk1",
						HostResource:    "/disk/vmdisk1",
						InstanceID:      6,
						Parent:          &tre,
						ResourceType:    17,
					},
				},
			},
		},
	}
	buf, _ := ovf.Build()
	fmt.Println(string(buf))
}
