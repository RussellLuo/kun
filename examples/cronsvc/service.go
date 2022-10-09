package cronsvc

import (
	"log"
)

//go:generate kungen ./service.go Service

// Service is used for handling cron jobs.
type Service interface {
	//kun:cron expr='@every 5s'
	SendEmail()
}

type Handler struct{}

func (h *Handler) SendEmail() {
	log.Println("Sending an email")
}
