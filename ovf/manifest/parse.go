package manifest

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var hashLineRx = regexp.MustCompile(`^([^(]+)\s*\(([^)]+)\)\s*=\s*([a-fA-F0-9]+)\s*$`)

func ParseFile(filename string) (Content, error) {
	fcontent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Parse(fcontent)
}

func Parse(fileContent []byte) (Content, error) {
	var ret Content
	var fcontent = string(fileContent)

	for idx, row := range strings.Split(fcontent, "\n") {
		if strings.TrimSpace(row) == "" {
			continue
		} else if !hashLineRx.MatchString(row) {
			return ret, fmt.Errorf("manifest format error, line %d", idx)
		}

		values := hashLineRx.FindAllStringSubmatch(row, -1)
		if len(values) != 1 || len(values[0]) != 4 {
			return ret, fmt.Errorf("manifest format error, line %d", idx)
		}
		ret = append(ret, FileHash{
			Algorithm: strings.TrimSpace(values[0][1]),
			Name:      strings.TrimSpace(values[0][2]),
			Hash:      strings.TrimSpace(values[0][3]),
		})
	}

	return ret, nil
}
