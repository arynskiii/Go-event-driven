package main

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"os"
)

func main() {
	logger := watermill.NewStdLogger(false, false)
	cli := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: cli,
	}, logger)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			panic(err)
		}
	}()
	firstMessage := message.NewMessage(watermill.NewUUID(), []byte("50"))
	secondMessage := message.NewMessage(watermill.NewUUID(), []byte("100"))

	if err := publisher.Publish("progress", firstMessage, secondMessage); err != nil {
		panic(err)
	}
}
