package http

import (
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"tickets/worker"
)

type Handler struct {
	worker         worker.Worker
	RedisPublisher *redisstream.Publisher
}
