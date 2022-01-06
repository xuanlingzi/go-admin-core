package runtime

import (
	"github.com/xuanlingzi/go-admin-core/storage"
)

type FileStore struct {
	prefix string
	store  storage.AdapterFileStore
}

// String string输出
func (e *FileStore) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

// Upload Put file to fileStores
func (e FileStore) Upload(name, location string) (string, error) {
	return e.store.Upload(name, location)
}