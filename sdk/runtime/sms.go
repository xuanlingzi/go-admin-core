package runtime

import (
	"github.com/go-admin-team/go-admin-core/storage"
)

// NewSms 创建对应上下文缓存
func NewSms(prefix string, store storage.AdapterSms) storage.AdapterSms {
	return &Sms{
		prefix:          prefix,
		store:           store,
	}
}

type Sms struct {
	prefix          string
	store           storage.AdapterSms
}

// String string输出
func (e *Sms) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

// Send val by sms
func (e Sms) Send(phones []string, templateId string, params []string) error {
	return e.store.Send(phones, templateId, params)
}