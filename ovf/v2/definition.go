package v2

import "encoding/xml"

const (
	DiskFormatVMDKStreamOptimized = "http://www.vmware.com/interfaces/specifications/vmdk.html#streamOptimized"
)

type File struct {
	ID   string `xml:"ovf:id,attr"`
	Href string `xml:"ovf:href,attr"`
}

type References struct {
	File []File `xml:"File"`
}

type Disk struct {
	Capacity int64  `xml:"ovf:capacity,attr"`
	DiskID   string `xml:"ovf:diskId,attr"`
	FileRef  string `xml:"ovf:fileRef,attr"`
	Format   string `xml:"ovf:format,attr"`
}

type DiskSection struct {
	Info string `xml:"Info,omitempty"`
	Disk Disk   `xml:"Disk"`
}

type Network struct {
	Name        string `xml:"ovf:name,attr"`
	Description string `xml:"Description"`
}

type NetworkSection struct {
	Info    string  `xml:"Info,omitempty"`
	Network Network `xml:"Network"`
}

type ProductSection struct {
	Info       string `xml:"Info,omitempty"`
	Product    string `xml:"Product,omitempty"`
	Vendor     string `xml:"Vendor,omitempty"`
	Version    string `xml:"Version,omitempty"`
	ProductURL string `xml:"ProductUrl,omitempty"`
	VendorURL  string `xml:"VendorUrl,omitempty"`
}

type AnnotationSection struct {
	Info       string `xml:"Info,omitempty"`
	Annotation string `xml:"Annotation,omitempty"`
}

type EulaSection struct {
	Info    string `xml:"Info,omitempty"`
	License string `xml:"License,omitempty"`
}

type OperatingSystemSection struct {
	ID          string `xml:"ovf:id,attr"`
	Info        string `xml:"Info,omitempty"`
	Description string `xml:"Description"`
}

type System struct {
	ElementName             string `xml:"vssd:ElementName"`
	InstanceID              int    `xml:"vssd:InstanceID"`
	VirtualSystemIdentifier string `xml:"vssd:VirtualSystemIdentifier"`
	VirtualSystemType       string `xml:"vssd:VirtualSystemType"`
}

type Item struct {
	Address         *int   `xml:"rasd:Address,omitempty"`
	Caption         string `xml:"rasd:Caption"`
	Description     string `xml:"rasd:Description"`
	InstanceID      int    `xml:"rasd:InstanceID"`
	ResourceType    int    `xml:"rasd:ResourceType"`
	ResourceSubType string `xml:"rasd:ResourceSubType,omitempty"`
	VirtualQuantity int    `xml:"rasd:VirtualQuantity"`
}

type StorageItem struct {
	AddressOnParent     string `xml:"sasd:AddressOnParent,omitempty"`
	Caption             string `xml:"sasd:Caption"`
	Description         string `xml:"sasd:Description"`
	HostResource        string `xml:"sasd:HostResource,omitempty"`
	InstanceID          int    `xml:"sasd:InstanceID"`
	Parent              *int   `xml:"sasd:Parent,omitempty"`
	ResourceType        int    `xml:"sasd:ResourceType"`
	AutomaticAllocation bool   `xml:"sasd:AutomaticAllocation,omitempty"`
}

type EthernetPortItem struct {
	AutomaticAllocation bool   `xml:"epasd:AutomaticAllocation,omitempty"`
	Caption             string `xml:"epasd:Caption"`
	Connection          string `xml:"epasd:Connection"`
	InstanceID          int    `xml:"epasd:InstanceID"`
	ResourceType        int    `xml:"epasd:ResourceType"`
	ResourceSubType     string `xml:"epasd:ResourceSubType,omitempty"`
}

type VirtualHardwareSection struct {
	Info              string             `xml:"Info,omitempty"`
	System            System             `xml:"System"`
	Items             []Item             `xml:"Item"`
	StorageItems      []StorageItem      `xml:"StorageItem"`
	EthernetPortItems []EthernetPortItem `xml:"EthernetPortItem"`
}

type VirtualSystem struct {
	ID                     string                 `xml:"ovf:id,attr"`
	Info                   string                 `xml:"Info,omitempty"`
	ProductSection         ProductSection         `xml:"ProductSection"`
	AnnotationSection      AnnotationSection      `xml:"AnnotationSection"`
	EulaSection            EulaSection            `xml:"EulaSection"`
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
	XMLEPASD   string `xml:"xmlns:epasd,attr"`

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
	c.XMLEPASD = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_StorageAllocationSettingData.xsd"

	if c.XMLLang == "" {
		c.XMLLang = "en-US"
	}
	return xml.Marshal(c)
}
