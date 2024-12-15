package ethereum

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client   *ethclient.Client
	poolAddr common.Address
}

type SwapEvent struct {
	Sender      common.Address
	Amount0In   *big.Int
	Amount1In   *big.Int
	Amount0Out  *big.Int
	Amount1Out  *big.Int
	BlockNumber uint64
	TxHash      common.Hash
}

func NewClient(config Config) (*Client, error) {
	client, err := ethclient.Dial(config.NodeURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:   client,
		poolAddr: common.HexToAddress(config.PoolAddress),
	}, nil
}

// 監聽 Swap 事件
func (c *Client) WatchSwaps(ctx context.Context, eventChan chan<- *SwapEvent) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{c.poolAddr},
		Topics: [][]common.Hash{{
			common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"), // Swap 事件的 topic
		}},
	}

	logs := make(chan types.Log)
	sub, err := c.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return err
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				log.Printf("訂閱錯誤: %v", err)
				return
			case vLog := <-logs:
				event := &SwapEvent{
					Sender:      common.HexToAddress(vLog.Topics[1].Hex()),
					Amount0In:   new(big.Int).SetBytes(vLog.Data[:32]),
					Amount1In:   new(big.Int).SetBytes(vLog.Data[32:64]),
					Amount0Out:  new(big.Int).SetBytes(vLog.Data[64:96]),
					Amount1Out:  new(big.Int).SetBytes(vLog.Data[96:128]),
					BlockNumber: vLog.BlockNumber,
					TxHash:      vLog.TxHash,
				}
				eventChan <- event
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
