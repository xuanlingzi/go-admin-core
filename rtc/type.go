package rtc

type AdapterZegoService interface {
	String() string
	Close() error
	Call(action, method string, query map[string]string, payload interface{}) (map[string]interface{}, error)
	Get(action string, query map[string]string) (map[string]interface{}, error)
	Post(action string, payload interface{}) (map[string]interface{}, error)
	GetClient() interface{}
}
