package moderation

const (
	PrefixKey = "__host"
)

type AdapterModeration interface {
	String() string
	AuditText(content string, suggestion *string, label *string, detail *string) error
	AuditImage(url string, suggestion *string, label *string, detail *string) error
	AuditVideo(url string, frame int32) error
	GetClient() interface{}
}
