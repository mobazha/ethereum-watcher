package ethereum_watcher

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/gommon/log"
	"github.com/mobazha/ethereum-watcher/plugin"
	"github.com/mobazha/ethereum-watcher/structs"
	"github.com/sirupsen/logrus"
)

// todo why some tx index in block is zero?
func TestTxReceiptPlugin(t *testing.T) {
	log.SetLevel(log.DEBUG)

	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)
	if err != nil {
		logrus.Panicln("RPC error:", err)
	}

	w.RegisterTxReceiptPlugin(plugin.NewTxReceiptPlugin(func(txAndReceipt *structs.RemovableTxAndReceipt) {
		if txAndReceipt.IsRemoved {
			fmt.Println("Removed >>", txAndReceipt.Tx.Hash(), txAndReceipt.Receipt.TransactionIndex)
		} else {
			fmt.Println("Adding >>", txAndReceipt.Tx.Hash(), txAndReceipt.Receipt.TransactionIndex)
		}
	}))

	w.RunTillExit()
}

func TestFilterPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)

	if err != nil {
		logrus.Panicln("RPC error:", err)
	}

	callback := func(txAndReceipt *structs.RemovableTxAndReceipt) {
		fmt.Println("tx:", txAndReceipt.Tx.Hash())
	}

	// only accept txs which end with: f
	filterFunc := func(tx *types.Transaction) bool {
		txHash := tx.Hash().Hex()

		return txHash[len(txHash)-1:] == "f"
	}

	w.RegisterTxReceiptPlugin(plugin.NewTxReceiptPluginWithFilter(callback, filterFunc))

	err = w.RunTillExitFromBlock(7840000)
	if err != nil {
		fmt.Println("RunTillExit with err:", err)
	}
}

func TestFilterPluginForDyDxApprove(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)

	if err != nil {
		logrus.Panicln("RPC error:", err)
	}

	callback := func(txAndReceipt *structs.RemovableTxAndReceipt) {
		receipt := txAndReceipt.Receipt

		for _, log := range receipt.Logs {
			topics := log.Topics
			if len(topics) == 3 &&
				topics[0].Hex() == "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925" &&
				topics[2].Hex() == "0x0000000000000000000000001e0447b19bb6ecfdae1e4ae1694b0c3659614e4e" {
				fmt.Printf(">> approving to dydx, tx: %s\n", txAndReceipt.Tx.Hash())
			}
		}
	}

	// only accept txs which send to DAI
	filterFunc := func(tx *types.Transaction) bool {
		return tx.To().Hex() == "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
	}

	w.RegisterTxReceiptPlugin(plugin.NewTxReceiptPluginWithFilter(callback, filterFunc))

	err = w.RunTillExitFromBlock(7844853)
	if err != nil {
		fmt.Println("RunTillExit with err:", err)
	}
}
