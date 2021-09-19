package follower

import (
	"context"
	"fmt"
	"strings"
	"tehranifar/fflow/storage"
	"tehranifar/fflow/tags"
	"time"

	"github.com/onflow/flow-go-sdk"
	flowClient "github.com/onflow/flow-go-sdk/client"
)

const blockTime = 3 * time.Second
const retryInterval = 1 * time.Second

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

		var currentBlock *flow.Block
		try(func() error {
			var err error
			currentBlock, err = f.client.GetLatestBlock(f.ctx, true)
			return err
		})

		currentHeight := currentBlock.Height

		for h := currentBlock.Height; h > height; h-- {
			fmt.Println(fmt.Sprintf("Processing block height %d", currentBlock.Height))
			try(func() error {
				return f.processBlock(f.ctx, currentBlock)
			})

			try(func() error {
				var err error
				currentBlock, err = f.client.GetBlockByHeight(f.ctx, h-1)
				return err
			})
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

			importTags := tags.ProcessImportTags(string(tx.Script))
			for _, tag := range importTags {
				failed := false
				if txRes.Error != nil {
					failed = true
				}
				RegisterTagMetrics(tag, failed)
			}
			dbTags := strings.Join(importTags, ",")

			f.storage.Save(&storage.Transaction{
				Tx:    txID.String(),
				Code:  string(tx.Script),
				Error: txError,
				Tags:  dbTags,
			})

		}
	}

	return nil
}

func try(f func() error) {
	err := f()
	for err != nil {
		fmt.Println(err)
		time.Sleep(retryInterval)
		err = f()
	}
}
