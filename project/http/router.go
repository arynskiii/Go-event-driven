package http

import (
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/labstack/echo/v4"
	"tickets/worker"
)

func NewHttpRouter(
	publisher *redisstream.Publisher,
	w *worker.Worker,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		RedisPublisher: publisher,
		worker:         *w,
	}

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
