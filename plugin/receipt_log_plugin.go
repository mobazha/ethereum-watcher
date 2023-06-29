package plugin

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mobazha/ethereum-watcher/structs"
)

type IReceiptLogPlugin interface {
	FromContracts() []common.Address
	InterestedTopics() [][]common.Hash
	NeedReceiptLog(receiptLog *structs.RemovableReceiptLog) bool
	Accept(receiptLog *structs.RemovableReceiptLog)
}

type ReceiptLogPlugin struct {
	contracts []common.Address
	topics    [][]common.Hash
	callback  func(receiptLog *structs.RemovableReceiptLog)
}

func NewReceiptLogPlugin(
	contracts []common.Address,
	topics [][]common.Hash,
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

func (p *ReceiptLogPlugin) InterestedTopics() [][]common.Hash {
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
	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Examples:
	// {} or nil          matches any topic list
	// {{A}}              matches topic A in first position
	// {{}, {B}}          matches any topic in first position AND B in second position
	// {{A}, {B}}         matches topic A in first position AND B in second position
	// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position

	if len(p.topics) == 0 {
		return false
	}

	if len(receiptLog.Topics) < len(p.topics) {
		return false
	}

	for pos, posTopics := range p.topics {
		posMatched := false
		if len(posTopics) == 0 {
			posMatched = true
		}

		for _, posTopic := range posTopics {
			if receiptLog.Topics[pos].Big().Cmp(posTopic.Big()) == 0 {
				posMatched = true
				break
			}
		}

		if !posMatched {
			return false
		}
	}

	return true
}
