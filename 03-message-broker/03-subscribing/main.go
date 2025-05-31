package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
	"os"
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	cli := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: cli,
	}, logger)
	if err != nil {
		logger.Error("Failed to initialize subscriber", err, nil)
		return
	}

	defer func() {
		if err := subscriber.Close(); err != nil {
			logger.Error("Failed to close subscriber", err, nil)
			return
		}
	}()

	messages, err := subscriber.Subscribe(context.Background(), "progress")
	for message := range messages {
		orderID := string(message.Payload)
		fmt.Printf("Message ID: %s - %s%", message.UUID, orderID)
		message.Ack()
	}
}
