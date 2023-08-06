package runtime

import (
	"github.com/xuanlingzi/go-admin-core/moderation"
)

type Moderation struct {
	prefix     string
	moderation moderation.AdapterModeration
}

// NewModeration 创建对应上下文缓存
func NewModeration(prefix string, moderation moderation.AdapterModeration) moderation.AdapterModeration {
	return &Moderation{
		prefix:     prefix,
		moderation: moderation,
	}
}

// String string输出
func (e *Moderation) String() string {
	if e.moderation == nil {
		return ""
	}
	return e.moderation.String()
}

// AuditText 文本审核
func (e *Moderation) AuditText(content string, result *int, label *string, score *int, detail *string) error {
	return e.moderation.AuditText(content, result, label, score, detail)
}

// AuditImage 图片审核
func (e *Moderation) AuditImage(url string, result *int, label *string, score *int, detail *string) error {
	return e.moderation.AuditImage(url, result, label, score, detail)
}

// AuditVideo 视频审核
func (e *Moderation) AuditVideo(url string, frame int32, jobId *string) error {
	return e.moderation.AuditVideo(url, frame, jobId)
}

// AuditResult 审核结果
func (e *Moderation) AuditResult(body *[]byte, result *int, label *string, score *int, detail *string, jobId *string) error {
	return e.moderation.AuditResult(body, result, label, score, detail, jobId)
}

func (e *Moderation) GetClient() interface{} {
	return e.moderation.GetClient()
}
