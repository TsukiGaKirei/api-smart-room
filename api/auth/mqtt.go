package auth

import (
	"api-smart-room/schema"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

//for publish
const (
	projectID = "delta-coil-393803"
	topicID   = "smart-room-pub-sub"
)

//
//PublishMessage zpublish message to microcontroller
func PublishMessage(ctx context.Context, data string) {
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("./delta-coil-393803-555a240284e0.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	topic := client.Topic(topicID)
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(data),
	})
	_, err = result.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
	fmt.Println("Message published successfully.")
}

//ReceiveMessage yang dikirim dari microcontroller apa aja?
//location id
//if there's a person
//room temp
func ReceiveMessage(ctx context.Context, client *pubsub.Client, subscriptionID string) {
	sub := client.Subscription(subscriptionID)
	var response schema.RoomUpdate
	// Receive messages in a loop (this is a simple example)
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", msg.Data)
		_, err := fmt.Sscanf(string(msg.Data), "%d %f %d", &response.RoomId, &response.Temperature, &response.PersonCount)
		if err != nil {
			fmt.Println("Error parsing:", err)
			return
		}
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
	}
}
