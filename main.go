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

const (
	maxMsgSize = 1024 * 1024 * 16
	mainnetURL = "access.mainnet.nodes.onflow.org:9000"

	postgresDSN = "host=db user=ff password=notsecure!notsecure!notsecure dbname=ff port=5432"
)

func main() {
	ctx := context.Background()

	storage, err := storage.NewPostgresStorage(postgresDSN)
	if err != nil {
		panic(err)
	}

	client, err := flowClient.New(
		mainnetURL,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
		grpc.WithInsecure())
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
