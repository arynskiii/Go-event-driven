package worker

import (
	"context"
)

type Task int

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

type Message struct {
	Task     Task
	TicketID string
}

type Worker struct {
	queue           chan Message
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
}

func NewWorker(spreadsheetsAPI SpreadsheetsAPI, receiptsService ReceiptsService) *Worker {
	return &Worker{
		queue:           make(chan Message, 100),
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
	}
}

func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		w.queue <- m
	}
}

func (w *Worker) Run(ctx context.Context) {
	for msg := range w.queue {
		switch msg.Task {
		case TaskIssueReceipt:
			if err := w.receiptsService.IssueReceipt(ctx, msg.TicketID); err != nil {
				w.Send(msg)
			}
		case TaskAppendToTracker:
			if err := w.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{msg.TicketID}); err != nil {
				w.Send(msg)
			}
		}
	}
}
