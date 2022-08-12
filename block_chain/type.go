package block_chain

const (
	PrefixKey = "__block_chain"
)

type AdapterBroker interface {
	String() string
	Send(content string, callback string) (string, error)
	Status(content string, hash string) (string, error)
}
