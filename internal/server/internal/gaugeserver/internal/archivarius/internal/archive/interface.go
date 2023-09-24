package archive

type Archive interface {
	Save([]byte) error
	Restore() ([]byte, error)
}
