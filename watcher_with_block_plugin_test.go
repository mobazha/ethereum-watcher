package ethereum_watcher

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"gitlab.com/eth-stack/ethereum-watcher/plugin"
	"gitlab.com/eth-stack/ethereum-watcher/structs"
)

func TestNewBlockNumPlugin(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)

	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)

	if err != nil {
		logrus.Panicln("RPC error:", err)
	}

	logrus.Println("waiting for new block...")
	w.RegisterBlockPlugin(plugin.NewBlockNumPlugin(func(i uint64, b bool) {
		logrus.Printf(">> found new block: %d, is removed: %t", i, b)
	}))

	w.RunTillExit()
}

func TestSimpleBlockPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)

	if err != nil {
		logrus.Panicln("RPC error:", err)
	}

	w.RegisterBlockPlugin(plugin.NewSimpleBlockPlugin(func(block *structs.RemovableBlock) {
		logrus.Infof(">> %+v", block.Block)
	}))

	w.RunTillExit()
}
