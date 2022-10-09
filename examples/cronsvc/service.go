package cronsvc

import (
	"context"
	"log"
)

//go:generate kungen ./service.go Service

// Service is used for handling cron jobs.
type Service interface {
	//kun:cron expr='@every 5s'
	SendEmail(ctx context.Context) error
}

type Handler struct{}

func (h *Handler) SendEmail(ctx context.Context) error {
	log.Println("Sending an email")
	return nil
}
