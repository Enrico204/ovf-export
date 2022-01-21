//nolint:forbidigo
package main

import (
	"flag"
	"fmt"
	"gitlab.com/enrico204/ovf-export/ovf/manifest"
	"gitlab.com/enrico204/ovf-export/ovf/ova"
	pu "gitlab.com/enrico204/ovf-export/proxmox-utils"
	"log"
	"os"
	"strconv"
)

func main() {
	var listVMs = flag.Bool("list", false, "List all VMs (same as qm list)")
	var vmIDparam = flag.String("id", "", "VM ID to export")
	var output = flag.String("output", "", "OVF destination directory")
	var cdrom = flag.Bool("cdrom", false, "Include CDROM/ISO images")
	var genOVA = flag.Bool("ova", false, "Generate OVA after OVF")
	var noManifest = flag.Bool("no-manifest", false, "Skip generating manifest")

	flag.Parse()

	switch {
	case !*listVMs && *output == "":
		log.Fatal("Output directory is needed")
	case *listVMs:
		// TODO: filter: only offline VMs
		out, err := pu.PrintListVMs()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		} else {
			fmt.Println(out)
		}
	case *vmIDparam != "":
		vmID, err := strconv.ParseInt(*vmIDparam, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		vm, err := pu.ExportVMInfo(vmID)
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
}
