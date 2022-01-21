package lu

import (
	"fmt"
	"github.com/digitalocean/go-libvirt"
)

func ExportDomain(l *libvirt.Libvirt, uuid libvirt.UUID) (Domain, error) {
	domain, err := l.DomainLookupByUUID(uuid)
	if err != nil {
		return Domain{}, fmt.Errorf("failed to lookup domain: %w", err)
	}

	domainxml, err := l.DomainGetXMLDesc(domain, libvirt.DomainXMLInactive)
	if err != nil {
		return Domain{}, fmt.Errorf("failed to get domain xml: %w", err)
	}

	domaindata, err := Parse([]byte(domainxml))
	if err != nil {
		return Domain{}, fmt.Errorf("failed to parse domain xml: %w", err)
	}

	return domaindata, nil
}
