package moderation

const (
	PrefixKey = "__host"
)

type AdapterModeration interface {
	String() string
	AuditText(content string, result *int, label *string, detail *string) error
	AuditImage(url string, result *int, label *string, detail *string) error
	AuditVideo(url string, frame int32, jobId *string) error
	AuditResult(body []byte, result *int, label *string, detail *string, jobId *string) error
	GetClient() interface{}
}
