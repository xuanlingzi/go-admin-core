package moderation

const (
	PrefixKey = "__host"
)

type AdapterModeration interface {
	String() string
	AuditText(content string, result *int, label *string, score *int, detail *string, jobId *string) error
	AuditImage(url string, fileSize int, result *int, label *string, score *int, detail *string, jobId *string) error
	AuditVideo(url string, frame int32, jobId *string) error
	AuditResult(body *[]byte, result *int, label *string, score *int, detail *string, jobId *string) error
	GetClient() interface{}
}
