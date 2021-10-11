package follower

import (
	"context"
	"fmt"
	"tehranifar/fflow/storage"
	"time"

	publisher "tehranifar/fflow/publisher"

	"cloud.google.com/go/pubsub"
	"github.com/onflow/flow-go-sdk"
	flowClient "github.com/onflow/flow-go-sdk/client"
)

const blockTime = 3 * time.Second
const retryInterval = 1 * time.Second

var (
	topic *pubsub.Topic
)

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
			fmt.Println("processign %s", txID)
			// txRes, err := f.client.GetTransactionResult(ctx, txID)
			// if err != nil {
			// 	return err
			// }

			// fmt.Println(fmt.Sprintf("processing txID res: %s !!!!!!!!", txRes))

			// tx, err := f.client.GetTransaction(ctx, txID)
			// if err != nil {
			// 	return err
			// }

			// txError := ""
			// failed := false
			// if txRes.Error != nil {
			// 	txError = txRes.Error.Error()
			// 	failed = true
			// }

			// var tagList []string
			// CadenceImports := GetImports(string(tx.Script))
			// for _, imp := range CadenceImports {
			// 	RegisterImportMetrics(imp, failed)
			// 	tagList = append(tagList, imp.Contract)
			// 	tagList = append(tagList, imp.Address)
			// }
			// dbImportTags := strings.Join(tagList, ",")

			// var authAddressList []string
			// for _, auth := range tx.Authorizers {
			// 	authAddressList = append(authAddressList, auth.String())
			// }
			// dbAuthorizer := strings.Join(authAddressList, ",")

			// var eventList []string
			// for _, event := range txRes.Events {
			// 	RegisterEventMetrics(event.Type)
			// 	eventList = append(eventList, event.Type)
			// }
			// dbEvents := strings.Join(eventList, ",")

			// f.storage.Save(&storage.Transaction{
			// 	Authorizers: dbAuthorizer,
			// 	Tx:          txID.String(),
			// 	Code:        string(tx.Script),
			// 	Error:       txError,
			// 	ImportTags:  dbImportTags,
			// 	Events:      dbEvents,
			// })
			publisher.Publish()

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
