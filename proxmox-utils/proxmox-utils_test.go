package pu

import "testing"

func TestParse(t *testing.T) {
	var vmcfg = `agent: 1
bootdisk: scsi0
cores: 2
ide2: none,media=cdrom
memory: 6144
name: test-env
net0: virtio=AA:BB:CC:11:22:33,bridge=vmbr0,firewall=1
numa: 0
onboot: 1
ostype: l26
scsi0: local-lvm:vm-100-disk-0,cache=writeback,size=200G
scsihw: virtio-scsi-pci
smbios1: uuid=82491b85-953f-4f9f-b817-717f174b1b1e
sockets: 2
vmgenid: 82491b85-953f-4f9f-b817-717f174b1b1e`

	vm, err := parseVMInfo(vmcfg)
	if err != nil {
		t.Fatal(err)
	}

	switch {
	case vm.Name != "test-env":
		t.Fatal("error reading vm name")
	case vm.Cores != 2:
		t.Fatal("error reading vm cores")
	case vm.Sockets != 2:
		t.Fatal("error reading vm sockets")
	case vm.SCSIHardware != "virtio-scsi-pci":
		t.Fatal("error reading scsi hardware")
	case vm.Memory != 6144:
		t.Fatal("error reading memory")
	case len(vm.Disks) != 1:
		t.Fatal("disk list length error")
	case vm.Disks[0].ID != 0:
		t.Fatal("SCSI disk #0 id error")
	case vm.Disks[0].Volume != "local-lvm:vm-100-disk-0":
		t.Fatal("SCSI disk #0 file name error")
	case vm.Disks[0].Cdrom:
		t.Fatal("SCSI disk #0 is cdrom error")
	case len(vm.Networks) != 1:
		t.Fatal("Networks list length error")
	case vm.Networks[0].ID != 0:
		t.Fatal("Networks #0 id error")
	case vm.Networks[0].Card != "virtio":
		t.Fatal("Networks #0 card error")
	}
}
