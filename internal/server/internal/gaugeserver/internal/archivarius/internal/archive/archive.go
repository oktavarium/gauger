package archive

import (
	"os"
)

type fileArchive struct {
	filename string
}

func NewFileArchive(filename string) Archive {
	return &fileArchive{filename: filename}
}

func (f *fileArchive) Save(data []byte) error {
	return os.WriteFile(f.filename, data, 0644)
}

func (f *fileArchive) Restore() ([]byte, error) {
	return os.ReadFile(f.filename)
}
