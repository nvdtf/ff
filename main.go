package main

import (
	"context"
	"fmt"

	flowClient "github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	flow, err := flowClient.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	block, err := flow.GetLatestBlock(ctx, true)
	if err != nil {
		panic(err)
	}

	for _, colGuarantee := range block.CollectionGuarantees {
		col, err := flow.GetCollection(ctx, colGuarantee.CollectionID)
		if err != nil {
			panic(err)
		}
		for _, txID := range col.TransactionIDs {
			tx, err := flow.GetTransaction(ctx, txID)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(tx.Script))
		}
	}

}
