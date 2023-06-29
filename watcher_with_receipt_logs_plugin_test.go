package ethereum_watcher

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mobazha/ethereum-watcher/plugin"
	"github.com/mobazha/ethereum-watcher/structs"
	"github.com/sirupsen/logrus"
)

type LogTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func TestReceiptLogsPlugin(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	api := "https://rpc.ankr.com/bsc_testnet_chapel"
	w, err := NewHttpBasedEthWatcher(context.Background(), api)
	if err != nil {
		logrus.Panicln(err)
	}

	contract := common.HexToAddress("0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee")
	topics := [][]common.Hash{{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}}

	ierc20Balance := `[{ "type" : "event", "name" : "Balance", "inputs": [{ "name" : "in", "type": "uint256" }] },
			{ "type" : "event", "name" : "Check", "inputs": [{ "name" : "t", "type": "address" }, { "name": "b", "type": "uint256" }] },
			{ "type" : "event", "name" : "Transfer", "inputs": [{ "name": "from", "type": "address", "indexed": true }, { "name": "to", "type": "address", "indexed": true }, { "name": "value", "type": "uint256" }] }
	]`
	contractAbi, err := abi.JSON(strings.NewReader(ierc20Balance))
	if err != nil {
		logrus.Panic(err)
	}

	w.RegisterReceiptLogPlugin(plugin.NewReceiptLogPlugin([]common.Address{contract}, topics, func(receipt *structs.RemovableReceiptLog) {
		if receipt.IsRemoved {
			logrus.Infof("Removed >> %+v", receipt)
		} else {
			logrus.Infof("Adding >> %+v, tx: %s, logIdx: %d %s", receipt, receipt.TxHash, receipt.Index, receipt.Data)
		}

		var log LogTransfer
		contractAbi.UnpackIntoInterface(&log, "Transfer", receipt.Data)

		log.From = common.HexToAddress(receipt.Topics[1].Hex())
		log.To = common.HexToAddress(receipt.Topics[2].Hex())

		logrus.Println("From", log.From)
		logrus.Println("To", log.To)
		logrus.Println("Value", log.Value)
	}))

	startBlock := 21881061
	err = w.RunTillExitFromBlock(uint64(startBlock))

	if err != nil {
		logrus.Panicln(err)
	}
}
