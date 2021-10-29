package runtime

import (
	"github.com/go-admin-team/go-admin-core/storage"
)

type Announce struct {
	prefix          string
	store           storage.AdapterAnnounce
}

// String string输出
func (e *Announce) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

// Send val by announces
func (e Announce) Send(addresses []string, template string, params []string) error {
	return e.store.Send(addresses, template, params)
}