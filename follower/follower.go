package follower

import (
	"context"
	"fmt"
	"tehranifar/fflow/storage"
	"time"

	"github.com/onflow/flow-go-sdk"
	flowClient "github.com/onflow/flow-go-sdk/client"
)

const blockTime = 3 * time.Second

type Follower struct {
	ctx     context.Context
	client  *flowClient.Client
	storage storage.Provider
}

func New(ctx context.Context, client *flowClient.Client, storage storage.Provider) Follower {
	return Follower{ctx, client, storage}
}

func (f *Follower) Follow(block *flow.Block) error {
	height := block.Height
	for {
		time.Sleep(blockTime)

		currentBlock, err := f.client.GetLatestBlock(f.ctx, true)
		if err != nil {
			return err
		}

		currentHeight := currentBlock.Height

		fmt.Println(fmt.Sprintf("Latest block height: %d", currentBlock.Height))

		for h := currentBlock.Height; h > height; h-- {
			fmt.Println(fmt.Sprintf("Processing block height %d", currentBlock.Height))
			err := f.processBlock(f.ctx, currentBlock)
			if err != nil {
				return err
			}

			currentBlock, err = f.client.GetBlockByHeight(f.ctx, h-1)
			if err != nil {
				return err
			}
		}

		height = currentHeight
	}

}

func (f *Follower) processBlock(ctx context.Context, block *flow.Block) error {
	for _, colGuarantee := range block.CollectionGuarantees {
		col, err := f.client.GetCollection(ctx, colGuarantee.CollectionID)
		if err != nil {
			return err
		}
		for _, txID := range col.TransactionIDs {
			txRes, err := f.client.GetTransactionResult(ctx, txID)
			if err != nil {
				return err
			}

			tx, err := f.client.GetTransaction(ctx, txID)
			if err != nil {
				return err
			}

			txError := ""
			if txRes.Error != nil {
				txError = txRes.Error.Error()
			}

			f.storage.Save(&storage.Transaction{
				Tx:    txID.String(),
				Code:  string(tx.Script),
				Error: txError,
			})

		}
	}

	return nil
}
