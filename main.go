package main

import (
	"context"
	"fmt"
	"os"
	"redis-tools/common"
	"redis-tools/stream"
)

func main() {
	ctx := context.Background()
	redisStream := stream.NewRedisStream(common.RDB, "my-stream", "my-consumer-group")
	if err := redisStream.InitConsumerGroup(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize consumer group: %v\n", err)
		os.Exit(1)
	}

	message := map[string]interface{}{"data": "hello world"}
	if _, err := redisStream.SendMessage(ctx, message); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send message: %v\n", err)
	}

	// Example to receive messages
	// This should be run in a separate process or go routine
	go redisStream.ReceiveMessages(ctx, "consumer1")

	// Prevent the main function from exiting immediately
	// This is just for demonstration; in a real application, you might have a more sophisticated mechanism
	select {}
}
