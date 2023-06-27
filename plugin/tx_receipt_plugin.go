package plugin

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mobazha/ethereum-watcher/structs"
)

type ITxReceiptPlugin interface {
	Accept(tx *structs.RemovableTxAndReceipt)
}

type TxReceiptPluginWithFilter struct {
	ITxReceiptPlugin
	filterFunc func(transaction *types.Transaction) bool
}

func (p TxReceiptPluginWithFilter) NeedReceipt(tx *types.Transaction) bool {
	return p.filterFunc(tx)
}

func NewTxReceiptPluginWithFilter(
	callback func(tx *structs.RemovableTxAndReceipt),
	filterFunc func(transaction *types.Transaction) bool) *TxReceiptPluginWithFilter {

	p := NewTxReceiptPlugin(callback)
	return &TxReceiptPluginWithFilter{p, filterFunc}
}

type TxReceiptPlugin struct {
	callback func(tx *structs.RemovableTxAndReceipt)
}

func NewTxReceiptPlugin(callback func(tx *structs.RemovableTxAndReceipt)) *TxReceiptPlugin {
	return &TxReceiptPlugin{callback}
}

func (p TxReceiptPlugin) Accept(tx *structs.RemovableTxAndReceipt) {
	if p.callback != nil {
		p.callback(tx)
	}
}
