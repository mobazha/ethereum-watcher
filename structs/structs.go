package structs

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type RemovableBlock struct {
	*types.Block
	IsRemoved bool
}

func NewRemovableBlock(block *types.Block, isRemoved bool) *RemovableBlock {
	return &RemovableBlock{
		block,
		isRemoved,
	}
}

type TxAndReceipt struct {
	Tx      *types.Transaction
	Receipt *types.Receipt
}

type RemovableTxAndReceipt struct {
	*TxAndReceipt
	IsRemoved bool
	TimeStamp uint64
}

type RemovableReceiptLog struct {
	*types.Log
	IsRemoved bool
}

func NewRemovableTxAndReceipt(tx *types.Transaction, receipt *types.Receipt, removed bool, timeStamp uint64) *RemovableTxAndReceipt {
	return &RemovableTxAndReceipt{
		&TxAndReceipt{
			tx,
			receipt,
		},
		removed,
		timeStamp,
	}
}

type RemovableTx struct {
	*types.Transaction
	IsRemoved bool
}

func NewRemovableTx(tx *types.Transaction, removed bool) RemovableTx {
	return RemovableTx{
		tx,
		removed,
	}
}

//
//type RemovableReceipt struct {
//	sdk.TransactionReceipt
//	IsRemoved bool
//}
//
//func NewRemovableReceipt(receipt sdk.TransactionReceipt, removed bool) RemovableReceipt {
//	return RemovableReceipt{
//		receipt,
//		removed,
//	}
//}
