package worker

import (
	"context"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
)

type Task int

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

type Message struct {
	Task     Task
	TicketID string
}

type Worker struct {
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
	subscriber      *redisstream.Subscriber
}

func NewWorker(spreadsheetsAPI SpreadsheetsAPI, receiptsService ReceiptsService, subscriber *redisstream.Subscriber) *Worker {
	return &Worker{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		subscriber:      subscriber,
	}
}

func (w *Worker) Run(ctx context.Context) {
	topics := []string{"issue-receipt", "append-to-tracker"}
	msgCh := make(chan *message.Message)

	for _, topic := range topics {
		go func(topic string) {
			for {
				messages, err := w.subscriber.Subscribe(ctx, topic)
				if err != nil {
					log.Fatalf("failed to subscribe to topic %s: %v", topic, err)
				}
				for {
					select {
					case <-ctx.Done():
						return
					case msg, ok := <-messages:
						if !ok {
							log.Printf("message channel closed for topic %s", topic)
							return
						}
						msg.Metadata.Set("topic", topic)
						msgCh <- msg
					}
				}
			}
		}(topic)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgCh:
			topic := msg.Metadata.Get("topic")
			switch topic {
			case "issue-receipt":
				if err := w.receiptsService.IssueReceipt(ctx, string(msg.Payload)); err != nil {
					msg.Nack()
					log.Printf("failed to issue receipt: %v", err)
				} else {
					msg.Ack()
				}
			case "append-to-tracker":
				if err := w.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{string(msg.Payload)}); err != nil {
					msg.Nack()
					log.Printf("failed to append row: %v", err)
				} else {
					msg.Ack()
				}
			}
		}
	}
}
