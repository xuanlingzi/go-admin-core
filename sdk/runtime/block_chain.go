package runtime

import "github.com/xuanlingzi/go-admin-core/block_chain"

type BlockChain struct {
	prefix string
	broker block_chain.AdapterBroker
}

// String string输出
func (e *BlockChain) String() string {
	if e.broker == nil {
		return ""
	}
	return e.broker.String()
}

// Send 发送上链内容
func (e BlockChain) Send(chain string, content string, callback string) (string, error) {
	return e.broker.Send(chain, content, callback)
}

// Status 上链状态
func (e BlockChain) Status(chain string, content string, hash string) (string, error) {
	return e.broker.Status(chain, content, hash)
}
