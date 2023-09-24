package archivarius

import "github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/archivarius/internal/storage"

type Archivarius interface {
	storage.Storage
}
