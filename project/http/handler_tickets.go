package http

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if err := h.RedisPublisher.Publish("issue-receipt", &message.Message{
			Payload: []byte(ticket),
		}); err != nil {
			return err
		}

		if err := h.RedisPublisher.Publish("append-to-tracker", &message.Message{
			Payload: []byte(ticket),
		}); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
