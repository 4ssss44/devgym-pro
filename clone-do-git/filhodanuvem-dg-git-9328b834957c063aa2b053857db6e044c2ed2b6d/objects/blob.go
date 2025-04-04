package objects

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/devgymbr/fit/hash"
)

func CreateBlob(f *os.File) (*os.File, error) {
	reader := bufio.NewReader(f)
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	h := hash.New(content)
	folder := h[0:2]
	fileName := h[2:]

	path := filepath.Join(".", ".fit", "objects", folder)
	if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(path, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	compressor := gzip.NewWriter(&b)
	if _, err := compressor.Write(content); err != nil {
		return nil, err
	}
	compressor.Close()
	_, err = io.WriteString(file, b.String())

	return file, err
}

func ReadBlob(hash string) ([]byte, error) {
	path := filepath.Join(".", ".fit", "objects", hash[0:2], hash[2:])
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return content, err
}
