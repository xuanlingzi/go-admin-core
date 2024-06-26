package runtime

import (
	"context"
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
func (e *FileStore) Upload(ctx context.Context, name, location string) (string, error) {
	return e.store.Upload(ctx, name, location)
}

// GetClient Put file to fileStores
func (e *FileStore) GetClient() interface{} {
	return e.store.GetClient()
}
