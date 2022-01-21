package pu

import (
	"os/exec"
)

func PrintListVMs() (string, error) {
	out, err := exec.Command("qm", "list").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
