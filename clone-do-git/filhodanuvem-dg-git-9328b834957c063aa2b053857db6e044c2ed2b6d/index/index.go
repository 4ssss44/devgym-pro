package index

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Index struct {
	Objects map[string]Entry `json:"objects"`
}

type Entry struct {
	Hash         string    `json:"hash"`
	Path         string    `json:"path"`
	LastModified time.Time `json:"last_modified"`
}

func Build() (Index, error) {
	var idx Index
	file, err := os.OpenFile(filepath.Join(".", ".fit", "index"), os.O_RDONLY, 0644)
	if err != nil {
		return idx, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return idx, err
	}

	err = json.Unmarshal(content, &idx)

	return idx, err
}

func Update(idx Index) error {
	content, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(".", ".fit", "index"), os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(content)

	return err
}

func AddBlob(file, blob *os.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	idx, err := Build()
	if err != nil {
		return err
	}

	parts := strings.Split(blob.Name(), "/")
	blobHash := parts[len(parts)-2] + parts[len(parts)-1]

	idx.Objects[file.Name()] = Entry{
		Hash:         blobHash,
		LastModified: fileInfo.ModTime(),
		Path:         file.Name(),
	}

	err = Update(idx)

	return err
}
