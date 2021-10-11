package publisher

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// func Publish() {

// 	// fmt.Println("Publishing message!!!!!!!")
// 	// ctx := context.Background()

// 	// msg := &pubsub.Message{
// 	// 	Data: []byte(r.FormValue("payload")),
// 	// }

// 	// if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
// 	// 	fmt.Sprintf("Could not publish message: %v", err)
// 	// 	return
// 	// }

// 	// fmt.Fprint(w, "Message published.")
// }

func Publish() {

	fmt.Println("Publishing message!!!!!!!")
	fmt.Print("******")

	projectID := "dapperlabs-data"
	topicID := "flow-events-poc"
	msg := "Hello World"

	ctx := context.Background()
	pubsub_client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Errorf("pubsub.NewClient: %v", err)
	} else {
		fmt.Println("created pubsub client: %s", pubsub_client)
	}
	defer pubsub_client.Close()

	t := pubsub_client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Errorf("Get: %v", err)
	}
	fmt.Println("Published a message; msg ID: %v\n", id)
	// return nil
}
