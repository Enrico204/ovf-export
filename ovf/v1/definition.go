package v1

import (
	"encoding/xml"
)

const (
	DiskFormatVMDKStreamOptimized = "http://www.vmware.com/interfaces/specifications/vmdk.html#streamOptimized"
	DiskFormatVMDKSparse          = "http://www.vmware.com/interfaces/specifications/vmdk.html#sparse"
	DiskFormatISO                 = "http://www.ecma-international.org/publications/standards/Ecma-119.htm"

	CompressionGzip     = "gzip"
	CompressionIdentity = "identity"
)

type File struct {
	Href        string `xml:"ovf:href,attr"`
	ID          string `xml:"ovf:id,attr"`
	Size        int64  `xml:"ovf:size,attr,omitempty"`
	Compression string `xml:"ovf:compression,omitempty,attr"`
	ChunkSize   int64  `xml:"ovf:chunkSize,omitempty,attr"`
}

type References struct {
	Files []File `xml:"File"`
}

type Disk struct {
	Capacity                int64  `xml:"ovf:capacity,attr"`
	CapacityAllocationUnits string `xml:"ovf:capacityAllocationUnits,omitempty,attr"`
	DiskID                  string `xml:"ovf:diskId,attr"`
	FileRef                 string `xml:"ovf:fileRef,attr"`
	Format                  string `xml:"ovf:format,attr"`
	PopulatedSize           int64  `xml:"ovf:populatedSize,omitempty,attr"`
}

type DiskSection struct {
	Info  string `xml:"Info"`
	Disks []Disk `xml:"Disk"`
}

type Network struct {
	Name string `xml:"ovf:name,attr"`

	Description string `xml:"Description"`
}

type NetworkSection struct {
	Info     string    `xml:"Info"`
	Networks []Network `xml:"Network"`
}

type ProductSection struct {
	Info       string `xml:"Info"`
	Product    string `xml:"Product,omitempty"`
	Vendor     string `xml:"Vendor,omitempty"`
	Version    string `xml:"Version,omitempty"`
	ProductURL string `xml:"ProductUrl,omitempty"`
	VendorURL  string `xml:"VendorUrl,omitempty"`
}

type AnnotationSection struct {
	Info       string `xml:"Info"`
	Annotation string `xml:"Annotation"`
}

type EulaSection struct {
	Info    string `xml:"Info"`
	License string `xml:"License"`
}

type OperatingSystemSection struct {
	ID string `xml:"ovf:id,attr"`

	Info        string `xml:"Info"`
	Description string `xml:"Description"`
}

type System struct {
	ElementName             string `xml:"vssd:ElementName"`
	InstanceID              int    `xml:"vssd:InstanceID"`
	VirtualSystemIdentifier string `xml:"vssd:VirtualSystemIdentifier"`
	VirtualSystemType       string `xml:"vssd:VirtualSystemType"`
}

type Item struct {
	Required *bool `xml:"ovf:required,omitempty"`

	AllocationUnits     string `xml:"rasd:AllocationUnits,omitempty"`
	AutomaticAllocation bool   `xml:"rasd:AutomaticAllocation,omitempty"`
	Address             *int   `xml:"rasd:Address,omitempty"`
	AddressOnParent     *int   `xml:"rasd:AddressOnParent,omitempty"`
	Caption             string `xml:"rasd:Caption"`
	Connection          string `xml:"rasd:Connection,omitempty"`
	Description         string `xml:"rasd:Description"`
	ElementName         string `xml:"rasd:ElementName"`
	HostResource        string `xml:"rasd:HostResource,omitempty"`
	InstanceID          int    `xml:"rasd:InstanceID"`
	Parent              *int   `xml:"rasd:Parent,omitempty"`
	ResourceSubType     string `xml:"rasd:ResourceSubType,omitempty"`
	ResourceType        int    `xml:"rasd:ResourceType"`
	VirtualQuantity     int64  `xml:"rasd:VirtualQuantity,omitempty"`
}

type VirtualHardwareSection struct {
	Info   string `xml:"Info,omitempty"`
	System System `xml:"System"`
	Items  []Item `xml:"Item"`
}

type VirtualSystem struct {
	ID                     string                 `xml:"ovf:id,attr"`
	Info                   string                 `xml:"Info"`
	ProductSection         *ProductSection        `xml:"ProductSection,omitempty"`
	AnnotationSection      *AnnotationSection     `xml:"AnnotationSection,omitempty"`
	EulaSections           []EulaSection          `xml:"EulaSection,omitempty"`
	OperatingSystemSection OperatingSystemSection `xml:"OperatingSystemSection"`
	VirtualHardwareSection VirtualHardwareSection `xml:"VirtualHardwareSection"`
}

type Envelope struct {
	XMLLang    string `xml:"xml:lang,attr"`
	OVFVersion string `xml:"ovf:version,attr"`
	XMLNS      string `xml:"xmlns,attr"`
	XMLNSOVF   string `xml:"xmlns:ovf,attr"`
	XMLNSRASD  string `xml:"xmlns:rasd,attr"`
	XMLNSVSSD  string `xml:"xmlns:vssd,attr"`
	XMLNSXSI   string `xml:"xmlns:xsi,attr"`

	References     References     `xml:"References"`
	DiskSection    DiskSection    `xml:"DiskSection"`
	NetworkSection NetworkSection `xml:"NetworkSection"`
	VirtualSystem  VirtualSystem  `xml:"VirtualSystem"`
}

func (c *Envelope) Build() ([]byte, error) {
	c.OVFVersion = "1.0"
	c.XMLNS = "http://schemas.dmtf.org/ovf/envelope/1"
	c.XMLNSOVF = "http://schemas.dmtf.org/ovf/envelope/1"
	c.XMLNSRASD = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData"
	c.XMLNSVSSD = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData"
	c.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"

	c.DiskSection.Info = "List of the virtual disks used in the package"
	c.NetworkSection.Info = "Logical networks used in the package"
	c.VirtualSystem.Info = "A virtual machine"
	if c.VirtualSystem.ProductSection != nil {
		c.VirtualSystem.ProductSection.Info = "Meta-information about the installed software"
	}
	if c.VirtualSystem.AnnotationSection != nil {
		c.VirtualSystem.AnnotationSection.Info = "A human-readable annotation"
	}
	for idx := range c.VirtualSystem.EulaSections {
		c.VirtualSystem.EulaSections[idx].Info = "License agreement for the virtual system"
	}
	c.VirtualSystem.OperatingSystemSection.Info = "The kind of installed guest operating system"
	c.VirtualSystem.VirtualHardwareSection.Info = "Virtual hardware requirements for a virtual machine"

	if c.XMLLang == "" {
		c.XMLLang = "en-US"
	}
	xmlbuf, err := xml.Marshal(c)
	if err != nil {
		return xmlbuf, err
	}

	return append([]byte(xml.Header), xmlbuf...), err
}
