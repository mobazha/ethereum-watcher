package ethereum_watcher

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"gitlab.com/eth-stack/ethereum-watcher/rpc"
)

type ReceiptLogWatcher struct {
	ctx                   context.Context
	api                   string
	startBlockNum         int
	contracts             []common.Address
	interestedTopics      []common.Hash
	handler               func(from, to int, receiptLogs []types.Log, isUpToHighestBlock bool) error
	config                ReceiptLogWatcherConfig
	highestSyncedBlockNum int
	highestSyncedLogIndex int
}

func NewReceiptLogWatcher(
	ctx context.Context,
	api string,
	startBlockNum int,
	contracts []common.Address,
	interestedTopics []common.Hash,
	handler func(from, to int, receiptLogs []types.Log, isUpToHighestBlock bool) error,
	configs ...ReceiptLogWatcherConfig,
) *ReceiptLogWatcher {

	config := decideConfig(configs...)

	pseudoSyncedLogIndex := config.StartSyncAfterLogIndex - 1

	return &ReceiptLogWatcher{
		ctx:                   ctx,
		api:                   api,
		startBlockNum:         startBlockNum,
		contracts:             contracts,
		interestedTopics:      interestedTopics,
		handler:               handler,
		config:                config,
		highestSyncedBlockNum: startBlockNum,
		highestSyncedLogIndex: pseudoSyncedLogIndex,
	}
}

func decideConfig(configs ...ReceiptLogWatcherConfig) ReceiptLogWatcherConfig {
	var config ReceiptLogWatcherConfig
	if len(configs) == 0 {
		config = defaultConfig
	} else {
		config = configs[0]

		if config.IntervalForPollingNewBlockInSec <= 0 {
			config.IntervalForPollingNewBlockInSec = defaultConfig.IntervalForPollingNewBlockInSec
		}

		if config.StepSizeForBigLag <= 0 {
			config.StepSizeForBigLag = defaultConfig.StepSizeForBigLag
		}

		if config.RPCMaxRetry <= 0 {
			config.RPCMaxRetry = defaultConfig.RPCMaxRetry
		}
	}

	return config
}

type ReceiptLogWatcherConfig struct {
	StepSizeForBigLag               int
	ReturnForBlockWithNoReceiptLog  bool
	IntervalForPollingNewBlockInSec int
	RPCMaxRetry                     int
	LagToHighestBlock               int
	StartSyncAfterLogIndex          int
}

var defaultConfig = ReceiptLogWatcherConfig{
	StepSizeForBigLag:               50,
	ReturnForBlockWithNoReceiptLog:  false,
	IntervalForPollingNewBlockInSec: 15,
	RPCMaxRetry:                     5,
	LagToHighestBlock:               0,
	StartSyncAfterLogIndex:          0,
}

func (w *ReceiptLogWatcher) Run() error {

	var blockNumToBeProcessedNext = w.startBlockNum

	rpc, err := rpc.NewEthRPCWithRetry(w.api, w.config.RPCMaxRetry)
	if err != nil {
		return err
	}

	for {
		select {
		case <-w.ctx.Done():
			return nil
		default:
			highestBlock, err := rpc.BlockNumber(w.ctx)
			if err != nil {
				return err
			}

			if blockNumToBeProcessedNext < 0 {
				blockNumToBeProcessedNext = int(highestBlock)
			}

			// [blockNumToBeProcessedNext...highestBlockCanProcess..[Lag]..CurrentHighestBlock]
			highestBlockCanProcess := int(highestBlock) - w.config.LagToHighestBlock
			numOfBlocksToProcess := highestBlockCanProcess - blockNumToBeProcessedNext + 1

			if numOfBlocksToProcess <= 0 {
				sleepSec := w.config.IntervalForPollingNewBlockInSec

				logrus.Debugf("no ready block after %d(lag: %d), sleep %d seconds", highestBlockCanProcess, w.config.LagToHighestBlock, sleepSec)

				select {
				case <-time.After(time.Duration(sleepSec) * time.Second):
					continue
				case <-w.ctx.Done():
					return nil
				}
			}

			var to int
			if numOfBlocksToProcess > w.config.StepSizeForBigLag {
				// quick mode
				to = blockNumToBeProcessedNext + w.config.StepSizeForBigLag - 1
			} else {
				// normal mode, up to cur highest block num can process
				to = highestBlockCanProcess
			}

			logs, err := rpc.FilterLogs(
				w.ctx,
				ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(blockNumToBeProcessedNext)),
					ToBlock:   big.NewInt(int64(to)),
					Addresses: w.contracts,
					Topics:    [][]common.Hash{w.interestedTopics},
				},
			)

			if err != nil {
				return err
			}

			isUpToHighestBlock := to == int(highestBlock)

			if len(logs) == 0 {
				if w.config.ReturnForBlockWithNoReceiptLog {
					err := w.handler(blockNumToBeProcessedNext, to, nil, isUpToHighestBlock)
					if err != nil {
						logrus.Infof("err when handling nil receipt log, block range: %d - %d", blockNumToBeProcessedNext, to)
						return fmt.Errorf("ethereum_watcher handler(nil) returns error: %s", err)
					}
				}
			} else {

				err := w.handler(blockNumToBeProcessedNext, to, logs, isUpToHighestBlock)
				if err != nil {
					logrus.Infof("err when handling receipt log, block range: %d - %d, receipt logs: %+v",
						blockNumToBeProcessedNext, to, logs,
					)

					return fmt.Errorf("ethereum_watcher handler returns error: %s", err)
				}
			}

			// todo rm 2nd param
			w.updateHighestSyncedBlockNumAndLogIndex(to, -1)

			blockNumToBeProcessedNext = to + 1
		}
	}
}

var progressLock = sync.Mutex{}

func (w *ReceiptLogWatcher) updateHighestSyncedBlockNumAndLogIndex(block int, logIndex int) {
	progressLock.Lock()
	defer progressLock.Unlock()

	w.highestSyncedBlockNum = block
	w.highestSyncedLogIndex = logIndex
}

func (w *ReceiptLogWatcher) GetHighestSyncedBlockNum() int {
	return w.highestSyncedBlockNum
}

func (w *ReceiptLogWatcher) GetHighestSyncedBlockNumAndLogIndex() (int, int) {
	progressLock.Lock()
	defer progressLock.Unlock()

	return w.highestSyncedBlockNum, w.highestSyncedLogIndex
}
