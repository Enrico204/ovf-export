package rasd

const (
	ResourceTypeOther                 = 1
	ResourceTypeOS                    = 2
	ResourceTypeCPU                   = 3
	ResourceTypeRAM                   = 4
	ResourceTypeIDE                   = 5
	ResourceTypeSCSI                  = 6
	ResourceTypeFC                    = 7
	ResourceTypeISCSI                 = 8
	ResourceTypeIBHCA                 = 9
	ResourceTypeEthernet              = 10
	ResourceTypeOtherEthernet         = 11
	ResourceTypeIOSlot                = 12
	ResourceTypeIODevice              = 13
	ResourceTypeFloppy                = 14
	ResourceTypeCDROM                 = 15
	ResourceTypeDVDROM                = 16
	ResourceTypeDisk                  = 17
	ResourceTypeTape                  = 18
	ResourceTypeStorageExtent         = 19
	ResourceTypeOtherStorageDevice    = 20
	ResourceTypeSerialPort            = 21
	ResourceTypeParallelPort          = 22
	ResourceTypeUSBController         = 23
	ResourceTypeGraphicsController    = 24
	ResourceTypeIEEE1394              = 25
	ResourceTypePartitionableUnit     = 26
	ResourceTypeBasePartitionableUnit = 27
	ResourceTypePowerSupply           = 28
	ResourceTypeCoolingDevice         = 29
	ResourceTypeEthernetSwitch        = 30
	ResourceTypeLogicalDisk           = 31
	ResourceTypeStorageVolume         = 32
	ResourceTypeEthernetConnection    = 33

	ResourceSubTypeIDEPIIX3    = "PIIX3"
	ResourceSubTypeIDEPIIX4    = "PIIX4"
	ResourceSubTypeLSILogic    = "lsilogic"
	ResourceSubTypeBusLogic    = "buslogic"
	ResourceSubTypeLSILogicSAS = "lsilogicsas"
	ResourceSubTypeSATA        = "AHCI"

	ResourceSubTypePCNet32 = "PCNet32"
	ResourceSubTypeE1000   = "E1000"
	ResourceSubTypeVMXNET3 = "VMXNET3"
)
