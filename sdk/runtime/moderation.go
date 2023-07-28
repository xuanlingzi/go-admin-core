package runtime

import (
	"github.com/xuanlingzi/go-admin-core/moderation"
)

type Moderation struct {
	prefix     string
	moderation moderation.AdapterModeration
}

// String string输出
func (e *Moderation) String() string {
	if e.moderation == nil {
		return ""
	}
	return e.moderation.String()
}

// AuditText 文本审核
func (e *Moderation) AuditText(content string, suggestion *string, label *string, detail *string) error {
	return e.moderation.AuditText(content, suggestion, label, detail)
}

// AuditImage 文本审核
func (e *Moderation) AuditImage(url string, suggestion *string, label *string, detail *string) error {
	return e.moderation.AuditImage(url, suggestion, label, detail)
}

// AuditVideo 文本审核
func (e *Moderation) AuditVideo(url string, frame int32) error {
	return e.moderation.AuditVideo(url, frame)
}

func (e *Moderation) GetClient() interface{} {
	return e.moderation.GetClient()
}
