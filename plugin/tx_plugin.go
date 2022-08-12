package plugin

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/eth-stack/ethereum-watcher/structs"
)

type ITxPlugin interface {
	AcceptTx(transaction structs.RemovableTx)
}

type TxHashPlugin struct {
	callback func(txHash common.Hash, isRemoved bool)
}

func (p TxHashPlugin) AcceptTx(transaction structs.RemovableTx) {
	if p.callback != nil {
		p.callback(transaction.Hash(), transaction.IsRemoved)
	}
}

func NewTxHashPlugin(callback func(txHash common.Hash, isRemoved bool)) TxHashPlugin {
	return TxHashPlugin{
		callback: callback,
	}
}

type TxPlugin struct {
	callback func(tx structs.RemovableTx)
}

func (p TxPlugin) AcceptTx(transaction structs.RemovableTx) {
	if p.callback != nil {
		p.callback(transaction)
	}
}

func NewTxPlugin(callback func(tx structs.RemovableTx)) TxPlugin {
	return TxPlugin{
		callback: callback,
	}
}
