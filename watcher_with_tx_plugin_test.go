package ethereum_watcher

import (
	"context"
	"fmt"
	"testing"

	"github.com/eth-stack/ethereum-watcher/plugin"
	"github.com/eth-stack/ethereum-watcher/structs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

func TestTxHashPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)

	if err != nil {
		logrus.Panicln("RPC error:", err)
	}

	w.RegisterTxPlugin(plugin.NewTxHashPlugin(func(txHash common.Hash, isRemoved bool) {
		fmt.Println(">>", txHash, isRemoved)
	}))

	w.RunTillExit()
}

func TestTxPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)

	if err != nil {
		logrus.Panicln("RPC error:", err)
	}

	w.RegisterTxPlugin(plugin.NewTxPlugin(func(tx structs.RemovableTx) {
		logrus.Printf(">>txHash: %s", tx.Hash())
	}))

	w.RunTillExit()
}
