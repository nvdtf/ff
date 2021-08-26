package main

import (
	"context"
	"net/http"
	"tehranifar/fflow/follower"

	flowClient "github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO LIST
// ---------
//
// - store transactions with gorm
// - test run locally with sqlite
// - deploy
//
// - store imports
//
// - figure out proper metrics
// - set up grafana to show metrics

func main() {
	ctx := context.Background()

	client, err := flowClient.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	startBlock, err := client.GetLatestBlock(ctx, true)
	if err != nil {
		panic(err)
	}

	f := follower.New(ctx, client)
	go f.Follow(startBlock)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
