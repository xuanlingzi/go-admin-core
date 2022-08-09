package block_chain

const (
	PrefixKey = "__block_chain"
)

type AdapterBroker interface {
	String() string
	Send(content string, callback string) error
	Status(content string, hash string) error
}
