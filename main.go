package main

import (
	"context"
	"fmt"
	"net/http"

	flowClient "github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO LIST
// ---------
//
// 1. write the transactions to a database with gorm -> write to grafana prom
// 2. infinite run
// 3. explore tools to analyze and graph the results
// 4. start with string checking a fixed list of known contract addresses

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

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
