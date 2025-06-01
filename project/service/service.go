package service

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/labstack/echo/v4"
	stdHTTP "net/http"
	"tickets/worker"

	ticketsHttp "tickets/http"
)

type Service struct {
	echoRouter *echo.Echo
	worker     *worker.Worker
}

func New(
	publisher *redisstream.Publisher,
	subscriber *redisstream.Subscriber,
	spreadsheetsService worker.SpreadsheetsAPI,
	receiptsService worker.ReceiptsService,
) Service {
	w := worker.NewWorker(spreadsheetsService, receiptsService, subscriber)
	echoRouter := ticketsHttp.NewHttpRouter(publisher, w)

	return Service{
		echoRouter: echoRouter,
		worker:     w,
	}
}

func (s Service) Run(ctx context.Context) error {
	go s.worker.Run(ctx)
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
