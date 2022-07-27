package lu

import "encoding/xml"

const (
	InterfaceTypeNetwork = "network"
	InterfaceTypeDirect  = "direct"

	EmulatorQEMUAMD64 = "/usr/bin/qemu-system-x86_64"
)

type Memory struct {
	Unit  string `xml:"unit,attr"`
	Value int64  `xml:",chardata"`
}

func (m *Memory) ToMegaBytes() int64 {
	switch m.Unit {
	case "KiB":
		return m.Value / 1024
	case "MiB":
		return m.Value
	case "GiB":
		return m.Value * 1024
	default:
		return 0
	}
}

type VCPU struct {
	Placement string `xml:"placement,attr"`
	Number    int    `xml:",chardata"`
}

type OperatingSystem struct {
	Type struct {
		Arch    string `xml:"arch,attr"`
		Machine string `xml:"machine,attr"`
		Value   string `xml:",chardata"`
	} `xml:"type"`
	Boot struct {
		Dev string `xml:"dev,attr"`
	} `xml:"boot"`
}

type Metadata struct {
	LibOSInfo struct {
		OS struct {
			ID string `xml:"id,attr"`
		} `xml:"os"`
	} `xml:"libosinfo"`
}

type Disk struct {
	Type   string `xml:"type,attr"`
	Device string `xml:"device,attr"`

	Driver struct {
		Name string `xml:"name,attr"`
		Type string `xml:"type,attr"`
	} `xml:"driver"`
	Source struct {
		File string `xml:"file,attr"`
	} `xml:"source"`
	Target struct {
		Dev string `xml:"dev,attr"`
		Bus string `xml:"bus,attr"`
	} `xml:"target"`
	ReadOnly bool `xml:"readonly"`
	Address  struct {
		Type       string `xml:"type,attr"`
		Controller string `xml:"controller,attr"`
		Bus        string `xml:"bus,attr"`
		Target     string `xml:"target,attr"`
		Unit       string `xml:"unit,attr"`
	} `xml:"address"`
}

type ControllerAddress struct {
	Type       string `xml:"type,attr"`
	Domain     string `xml:"domain,attr"`
	Controller string `xml:"controller,attr"`
	Bus        string `xml:"bus,attr"`
	Port       string `xml:"port,attr"`
	Slot       string `xml:"slot,attr"`
	Function   string `xml:"function,attr"`
}

type Controller struct {
	Type  string `xml:"type,attr"`
	Index int    `xml:"index,attr"`
	Model string `xml:"model,attr"`
	Ports int    `xml:"ports,attr"`

	Address ControllerAddress `xml:"address"`
}

type Interface struct {
	Type string `xml:"type,attr"`

	MAC struct {
		Address string `xml:"address,attr"`
	} `xml:"mac"`
	Source struct {
		Network string `xml:"network,attr"`
	} `xml:"source"`
	Model struct {
		Type string `xml:"type,attr"`
	} `xml:"model"`
	Address ControllerAddress `xml:"address"`
}

type Video struct {
	Model struct {
		Type    string `xml:"type,attr"`
		VRAM    int64  `xml:"vram,attr"`
		Heads   int    `xml:"heads,attr"`
		Primary string `xml:"primary,attr"`
	} `xml:"model"`
	Address ControllerAddress `xml:"address"`
}

type Graphics struct {
	Type     string `xml:"type,attr"`
	AutoPort string `xml:"autoport,attr"`

	Listen struct {
		Type string `xml:"type,attr"`
	} `xml:"listen"`
}

type Input struct {
	Type string `xml:"type,attr"`
	Bus  string `xml:"bus,attr"`
}

type Channel struct {
	Type string `xml:"type,attr"`

	Target struct {
		Type string `xml:"type,attr"`
		Name string `xml:"vram,attr"`
	} `xml:"target"`
	Address ControllerAddress `xml:"address"`
}

type MemBalloon struct {
	Model string `xml:"model,attr"`

	Address ControllerAddress `xml:"address"`
}

type RNG struct {
	Model string `xml:"model,attr"`

	Backend struct {
		Model string `xml:"model,attr"`
		Path  string `xml:",chardata"`
	} `xml:"backend"`
	Address ControllerAddress `xml:"address"`
}

type Domain struct {
	Type string `xml:"type,attr"`

	Name          string          `xml:"name"`
	UUID          string          `xml:"uuid"`
	Metadata      Metadata        `xml:"metadata"`
	Memory        Memory          `xml:"memory"`
	CurrentMemory Memory          `xml:"currentMemory"`
	VCPUs         VCPU            `xml:"vcpu"`
	OS            OperatingSystem `xml:"os"`
	Devices       struct {
		Emulator    string       `xml:"emulator"`
		Disks       []Disk       `xml:"disk"`
		Controllers []Controller `xml:"controller"`
		Interfaces  []Interface  `xml:"interface"`
		Videos      []Video      `xml:"video"`
		/*Graphics    []Graphics   `xml:"graphics"`
		Inputs      []Input      `xml:"input"`
		Channels    []Channel    `xml:"channel"`
		MemBalloon  MemBalloon   `xml:"memballoon"`
		RNGs        []RNG        `xml:"rng"`*/
	} `xml:"devices"`
}

func Parse(domainxml []byte) (Domain, error) {
	var ret Domain
	err := xml.Unmarshal(domainxml, &ret)
	return ret, err
}
