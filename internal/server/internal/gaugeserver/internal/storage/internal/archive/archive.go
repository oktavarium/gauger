package archive

import (
	"os"
)

type FileArchive struct {
	filename string
}

func NewFileArchive(filename string) FileArchive {
	return FileArchive{filename: filename}
}

func (f FileArchive) Save(data []byte) error {
	return os.WriteFile(f.filename, data, 0644)
}

func (f FileArchive) Restore() ([]byte, error) {
	return os.ReadFile(f.filename)
}
