# ovf-export for libvirt

`ovf-export` generates OVF/OVA files from libvirt domains. Generated OVA/OVF
files can be imported in other hypervisors, like VirtualBox or VMware
Player/Workstation/ESXi.

Also, `ovf` package can be used for serializing custom OVF files.
De-serializing is not supported at the moment. See "known issues" for details.

**Beta software**: this code works for me, but it may not work in your system.
At the moment, the code is not tested with different environments. Feel free to
propose patches.

## Usage

Requirements:
* `qemu-img`
* `tar` (only when generating OVA files)

## Installation

```shell
$ go install gitlab.com/enrico204/ovf-export/cmd/ovf-export
```

If you want to use this tool as library:

```shell
# Inside the project dir
$ go get gitlab.com/enrico204/ovf-export
```

## Usage

```shell
$ ovf-export -list
UUID                                    Name
--------------------------------------------------------
0962a56e96144bdd99bd83418c3e425a        debian10
61290c324b744d9095c6023520ddf72d        openindiana
e7904362d30b4179977804d073f070b4        win10

$ ovf-export -id 0962a56e96144bdd99bd83418c3e425a -output ~/debian10ovf -ova
```

Libvirt URL can be specified via env variable `LIBVIRT_URL` or using the
command line flag `-libvirt-url`. It can be the path for the local socket, or
`hostname:port` for TCP connection. The default value is `/var/run/libvirt/libvirt-sock`.

Use `-help` to list all command line flags:

```shell
$ ovf-export -help
  -cdrom
        Include CDROM/ISO images
  -id string
        Domain ID to export
  -libvirt-url string
        Libvirt URL
  -list
        List all inactive libvirt-utils domains (a.k.a. VMs)
  -no-manifest
        Skip generating manifest
  -output string
        OVF destination directory
  -ova
        Generate OVA after OVF
```

## Known issues

* When exporting a domain with an IDE controller, VMware Workstation refuses to
import  the OVA/OVF file due to "mismatch hash" for the VMDK (even if the hash
matches). A workaround is remove the manifest before importing the image.
* De-serializing a OVF file is not supported, and probably won't be supported 
soon due a [limit in Go `encoding/xml` handling namespaced XML prefixes](https://github.com/golang/go/issues/9519).
* Currently, the program should be in the same machine as the `libvirt` daemon.

# LICENSE

This code is released under MIT license. See [LICENSE](LICENSE) for details.
