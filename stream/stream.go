package stream

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStream struct {
	client        *redis.Client
	streamName    string
	consumerGroup string
}

func NewRedisStream(client *redis.Client, streamName, consumerGroup string) *RedisStream {
	return &RedisStream{
		client:        client,
		streamName:    streamName,
		consumerGroup: consumerGroup,
	}
}

func (rs *RedisStream) InitConsumerGroup(ctx context.Context) error {
	err := rs.client.XGroupCreateMkStream(ctx, rs.streamName, rs.consumerGroup, "$").Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("error creating consumer group: %w", err)
	}
	return nil
}

func (rs *RedisStream) SendMessage(ctx context.Context, message map[string]interface{}) (string, error) {
	id, err := rs.client.XAdd(ctx, &redis.XAddArgs{
		Stream: rs.streamName,
		Values: message,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("error adding message to stream: %w", err)
	}
	return id, nil
}

func (rs *RedisStream) ReceiveMessages(ctx context.Context, consumerName string) {
	for {
		streams, err := rs.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    rs.consumerGroup,
			Consumer: consumerName,
			Streams:  []string{rs.streamName, ">"},
			Block:    0,
		}).Result()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stream: %v\n", err)
			time.Sleep(1 * time.Second) // Wait a bit before trying again
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				fmt.Printf("Received message: %v\n", message.Values)
				// Acknowledge the message so it won't be read again
				_, err := rs.client.XAck(ctx, rs.streamName, rs.consumerGroup, message.ID).Result()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error acknowledging message: %v\n", err)
				}
			}
		}
	}
}
