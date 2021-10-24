package runtime

import (
	"github.com/go-admin-team/go-admin-core/storage"
)

// NewCos 创建对应上下文缓存
func NewCos(prefix string, store storage.AdapterCos) storage.AdapterCos {
	return &Cos{
		prefix:          prefix,
		store:           store,
	}
}

type Cos struct {
	prefix          string
	store           storage.AdapterCos
}

// String string输出
func (e *Cos) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

// Put file to cos
func (e Cos) PutFromFile(fileLocation string) error {
	return e.store.PutFromFile(fileLocation)
}