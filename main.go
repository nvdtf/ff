package main

import (
	"context"
	"net/http"
	"tehranifar/fflow/follower"
	"tehranifar/fflow/storage"

	flowClient "github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO LIST
// ---------
//
//
// - deploy to ec2
// - more contracts
// - create panels
// - event listener (from tx list)
// - event metrics

func main() {
	ctx := context.Background()

	storage, err := storage.NewSqliteStorage()
	if err != nil {
		panic(err)
	}

	client, err := flowClient.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	startBlock, err := client.GetLatestBlock(ctx, true)
	if err != nil {
		panic(err)
	}

	f := follower.New(ctx, client, storage)
	go f.Follow(startBlock)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
