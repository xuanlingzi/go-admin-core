package block_chain

const (
	PrefixKey = "__block_chain"
)

type AdapterBroker interface {
	String() string
	Send(chain string, content string, callback string) (string, error)
	Status(chain string, content string, hash string) (string, error)
}
