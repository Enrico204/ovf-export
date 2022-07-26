package lu

import "encoding/xml"

const (
	InterfaceTypeNetwork = "network"
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
	Type     string `xml:"type,attr"`
	Domain   string `xml:"domain,attr"`
	Bus      string `xml:"bus,attr"`
	Slot     string `xml:"slot,attr"`
	Function string `xml:"function,attr"`
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
		Disks       []Disk       `xml:"disk"`
		Controllers []Controller `xml:"controller"`
		Interfaces  []Interface  `xml:"interface"`
	} `xml:"devices"`
}

func Parse(domainxml []byte) (Domain, error) {
	var ret Domain
	err := xml.Unmarshal(domainxml, &ret)
	return ret, err
}
