//nolint:forbidigo
package main

import (
	"flag"
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	lu "gitlab.com/enrico204/ovf-export/libvirt-utils"
	"gitlab.com/enrico204/ovf-export/ovf/manifest"
	"gitlab.com/enrico204/ovf-export/ovf/ova"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var l *libvirt.Libvirt

func main() {
	var listDomains = flag.Bool("list", false, "List all inactive libvirt-utils domains (a.k.a. VMs)")
	var domainID = flag.String("id", "", "Domain ID to export")
	var output = flag.String("output", "", "OVF destination directory")
	var cdrom = flag.Bool("cdrom", false, "Include CDROM/ISO images")
	var genOVA = flag.Bool("ova", false, "Generate OVA after OVF")
	var noManifest = flag.Bool("no-manifest", false, "Skip generating manifest")
	var libvirtURL = flag.String("libvirt-url", os.Getenv("LIBVIRT_URL"), "Libvirt URL")

	flag.Parse()

	if !*listDomains && *output == "" {
		log.Fatal("Output directory is needed")
	}

	var libvirtNetwork = "tcp"

	if *libvirtURL == "" {
		*libvirtURL = "/var/run/libvirt/libvirt-sock"
	}
	if strings.Contains(*libvirtURL, "/") {
		libvirtNetwork = "unix"
	}

	// This dials libvirt-utils on the local machine, but you can substitute the first
	// two parameters with "tcp", "<ip address>:<port>" to connect to libvirt-utils on
	// a remote machine.
	c, err := net.DialTimeout(libvirtNetwork, *libvirtURL, 5*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt-utils: %v", err)
	}

	l = libvirt.NewWithDialer(dialers.NewAlreadyConnected(c))
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	if *listDomains {
		domains, _, err := l.ConnectListAllDomains(1, libvirt.ConnectListDomainsInactive)
		if err != nil {
			log.Fatalf("failed to retrieve domains: %v", err)
		}

		fmt.Println("UUID\t\t\t\t\tName")
		fmt.Printf("--------------------------------------------------------\n")
		for _, d := range domains {
			fmt.Printf("%x\t%s\n", d.UUID, d.Name)
		}
	} else if *domainID != "" {
		uuid, err := lu.DecodeLibvirtUUID(*domainID)
		if err != nil {
			log.Fatal(err)
		}

		vm, err := lu.ExportDomain(l, uuid)
		if err != nil {
			log.Fatal(err)
		}

		if err := encodeToOvf(vm, *cdrom, *output); err != nil {
			log.Fatal(err)
		}

		if !*noManifest {
			if err := manifest.GenerateManifestFromOVFDir(*output); err != nil {
				log.Fatal(err)
			}
		}

		if *genOVA {
			if err := ova.GenerateFromOVFDir(*output); err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := l.Disconnect(); err != nil {
		log.Fatalf("failed to disconnect: %v", err)
	}
}
