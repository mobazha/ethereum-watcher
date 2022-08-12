package plugin

import (
	"github.com/eth-stack/ethereum-watcher/structs"
	"github.com/ethereum/go-ethereum/common"
)

type IReceiptLogPlugin interface {
	FromContracts() []common.Address
	InterestedTopics() []common.Hash
	NeedReceiptLog(receiptLog *structs.RemovableReceiptLog) bool
	Accept(receiptLog *structs.RemovableReceiptLog)
}

type ReceiptLogPlugin struct {
	contracts []common.Address
	topics    []common.Hash
	callback  func(receiptLog *structs.RemovableReceiptLog)
}

func NewReceiptLogPlugin(
	contracts []common.Address,
	topics []common.Hash,
	callback func(receiptLog *structs.RemovableReceiptLog),
) *ReceiptLogPlugin {
	return &ReceiptLogPlugin{
		contracts: contracts,
		topics:    topics,
		callback:  callback,
	}
}

func (p *ReceiptLogPlugin) FromContracts() []common.Address {
	return p.contracts
}

func (p *ReceiptLogPlugin) InterestedTopics() []common.Hash {
	return p.topics
}

func (p *ReceiptLogPlugin) Accept(receiptLog *structs.RemovableReceiptLog) {
	if p.callback != nil {
		p.callback(receiptLog)
	}
}

// simplified version of specifying topic filters
// https://github.com/ethereum/wiki/wiki/JSON-RPC#a-note-on-specifying-topic-filters
func (p *ReceiptLogPlugin) NeedReceiptLog(receiptLog *structs.RemovableReceiptLog) bool {
	var firstTopic common.Hash
	if len(receiptLog.Topics) > 0 {
		firstTopic = receiptLog.Topics[0]
	}

	for _, interestedTopic := range p.topics {
		if firstTopic.Big().Cmp(interestedTopic.Big()) == 0 {
			return true
		}
	}

	return false
}
