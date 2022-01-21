package main

import (
	"fmt"
	"gitlab.com/enrico204/ovf-export/ovf/rasd"
	v1 "gitlab.com/enrico204/ovf-export/ovf/v1"
	"gitlab.com/enrico204/ovf-export/ovf/vssd"
	pu "gitlab.com/enrico204/ovf-export/proxmox-utils"
	"io/ioutil"
	"path"
	"strconv"
)

func encodeToOvf(vm pu.VM, cdrom bool, outdir string) error {
	var zero = 0
	var vcpus = vm.Sockets * vm.Cores

	var items = []v1.Item{
		{
			Caption:         fmt.Sprintf("%d virtual CPU", vcpus),
			Description:     "Number of virtual CPUs",
			ElementName:     fmt.Sprintf("%d virtual CPU", vcpus),
			InstanceID:      1,
			ResourceType:    rasd.ResourceTypeCPU,
			VirtualQuantity: vcpus,
		},
		{
			AllocationUnits: "MegaBytes",
			Caption:         fmt.Sprintf("%d MB of memory", vm.Memory),
			Description:     "Memory Size",
			ElementName:     fmt.Sprintf("%d MB of memory", vm.Memory),
			InstanceID:      2,
			ResourceType:    rasd.ResourceTypeRAM,
			VirtualQuantity: vm.Memory,
		},
	}
	var nextInstanceID = 3
	var nextControllerAddress = 0

	var sataController *v1.Item
	var sataNextAddress = 0

	var referencedFiles []v1.File
	var disks []v1.Disk
	for idx, disk := range vm.Disks {
		if !disk.Cdrom || (disk.Cdrom && cdrom) {
			fileRef := fmt.Sprintf("disk%d", idx)

			sourceFile, err := disk.GetFilePath()
			if err != nil {
				return err
			}

			convStat, err := importDisk(sourceFile, outdir, disk.Bus)
			if err != nil {
				return err
			}

			referencedFiles = append(referencedFiles, v1.File{
				ID:   fileRef,
				Href: path.Base(convStat.FileName),
				Size: convStat.FileSize,
			})

			var diskFormat = v1.DiskFormatVMDKStreamOptimized
			var diskType = rasd.ResourceTypeDisk
			if disk.Cdrom {
				diskFormat = v1.DiskFormatISO
				diskType = rasd.ResourceTypeCDROM
			}

			manifestDisk := v1.Disk{
				DiskID:  fileRef,
				FileRef: fileRef,
				Format:  diskFormat,
			}
			if convStat.TotalCapacity > 0 && convStat.SparseFileOccupation > 0 {
				manifestDisk.PopulatedSize = convStat.SparseFileOccupation
				manifestDisk.Capacity = convStat.TotalCapacity
			} else {
				manifestDisk.Capacity = convStat.FileSize
			}
			disks = append(disks, manifestDisk)

			if disk.Bus == pu.DiskBusSCSI {
				if sataController == nil {
					controllerAddress := nextControllerAddress
					sataController = &v1.Item{
						Address:         &controllerAddress,
						Caption:         "sataController0",
						Description:     "SATA Controller",
						ElementName:     "sataController0",
						InstanceID:      nextInstanceID,
						ResourceType:    rasd.ResourceTypeSCSI,
						ResourceSubType: rasd.ResourceSubTypeSATA,
					}
					nextInstanceID++
					items = append(items, *sataController)
				}

				sataAddress := sataNextAddress
				items = append(items, v1.Item{
					Parent:          &sataController.InstanceID,
					AddressOnParent: &sataAddress,
					Caption:         fileRef,
					ElementName:     fileRef,
					HostResource:    "/disk/" + fileRef,
					InstanceID:      nextInstanceID,
					ResourceType:    diskType,
				})

				sataNextAddress++
				nextInstanceID++
			} else if disk.Bus == pu.DiskBusIDE {
				controllerInstanceID := nextInstanceID
				controllerAddress := nextControllerAddress

				items = append(items, v1.Item{
					Address:         &controllerAddress,
					Caption:         "ideController0",
					Description:     "IDE Controller",
					ElementName:     "ideController0",
					InstanceID:      nextInstanceID,
					ResourceType:    rasd.ResourceTypeIDE,
					ResourceSubType: rasd.ResourceSubTypeIDEPIIX4,
				}, v1.Item{
					Parent:          &controllerInstanceID,
					AddressOnParent: &zero,
					Caption:         fileRef,
					ElementName:     fileRef,
					HostResource:    "/disk/" + fileRef,
					InstanceID:      nextInstanceID + 1,
					ResourceType:    diskType,
				})

				nextInstanceID += 2
				nextControllerAddress++
			}
		}
	}

	var networks = make([]v1.Network, 0, len(vm.Networks))
	for _, intf := range vm.Networks {
		var nic = "net" + strconv.FormatInt(intf.ID, 10)
		networks = append(networks, v1.Network{
			Name:        nic,
			Description: intf.Card,
		})

		var subtype = rasd.ResourceSubTypeE1000
		if intf.Card == pu.NetVMXNET3 {
			subtype = rasd.ResourceSubTypeVMXNET3
		}
		items = append(items, v1.Item{
			Caption:             "Ethernet NIC " + nic,
			Connection:          nic,
			ElementName:         "Ethernet NIC " + nic,
			InstanceID:          nextInstanceID,
			ResourceType:        rasd.ResourceTypeEthernet,
			ResourceSubType:     subtype,
			AutomaticAllocation: true,
		})

		nextInstanceID++
	}

	var ovf = v1.Envelope{
		XMLLang:     "en-US",
		References:  v1.References{Files: referencedFiles},
		DiskSection: v1.DiskSection{Disks: disks},
		NetworkSection: v1.NetworkSection{
			Networks: networks,
		},
		VirtualSystem: v1.VirtualSystem{
			ID: vm.Name,
			OperatingSystemSection: v1.OperatingSystemSection{
				ID:          "0",
				Description: "Other",
			},
			VirtualHardwareSection: v1.VirtualHardwareSection{
				System: v1.System{
					ElementName:             "Virtual Hardware Family",
					InstanceID:              0,
					VirtualSystemIdentifier: vm.Name,
					VirtualSystemType:       vssd.VirtualSystemIdentifierVMX07,
				},
				Items: items,
			},
		},
	}

	ovfbuf, err := ovf.Build()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(outdir, vm.Name+".ovf"), ovfbuf, 0600)
}
