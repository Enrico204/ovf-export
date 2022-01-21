package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const (
	AlgorithmSHA256 = "SHA256"
	AlgorithmSHA1   = "SHA1"
)

type FileHash struct {
	Algorithm string
	Name      string
	Hash      string
}

type Content []FileHash

func (c *Content) AddFile(filepath string) error {
	fname := path.Base(filepath)

	fp, err := os.Open(filepath)
	if err != nil {
		return err
	}

	var sum = sha256.New()
	_, err = io.Copy(sum, fp)
	if err != nil {
		return err
	}

	*c = append(*c, FileHash{
		Algorithm: AlgorithmSHA256,
		Name:      fname,
		Hash:      hex.EncodeToString(sum.Sum(nil)),
	})
	return nil
}

func (c *Content) Build() ([]byte, error) {
	var ret strings.Builder
	var err error
	for _, row := range *c {
		_, err = fmt.Fprintf(&ret, "%s(%s)= %s\n", row.Algorithm, row.Name, row.Hash)
		if err != nil {
			return nil, err
		}
	}
	return []byte(ret.String()), nil
}
