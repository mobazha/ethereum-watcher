package rpc

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthBlockChainRPCWithRetry struct {
	*ethclient.Client
	maxRetryTimes int
}

func NewEthRPCWithRetry(api string, maxRetryCount int) (*EthBlockChainRPCWithRetry, error) {
	rpc, err := ethclient.Dial(api)
	if err != nil {
		return nil, err
	}

	return &EthBlockChainRPCWithRetry{rpc, maxRetryCount}, nil
}

func (rpc EthBlockChainRPCWithRetry) BlockByNumber(ctx context.Context, number *big.Int) (rst *types.Block, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.Client.BlockByNumber(context.TODO(), number)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) HeaderByNumber(ctx context.Context, number *big.Int) (rst *types.Header, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.Client.HeaderByNumber(ctx, number)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) TransactionReceipt(ctx context.Context, txHash common.Hash) (rst *types.Receipt, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.Client.TransactionReceipt(ctx, txHash)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) BlockNumber(ctx context.Context) (rst uint64, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.Client.BlockNumber(ctx)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) FilterLogs(ctx context.Context, query ethereum.FilterQuery) (rst []types.Log, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.Client.FilterLogs(ctx, query)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}
