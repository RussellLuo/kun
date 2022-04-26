package eventsvc

import (
	"context"
	"fmt"
)

//go:generate kungen ./service.go Service

// Service is used for handling events.
type Service interface {
	//kun:event type=created
	EventCreated(ctx context.Context, id int) error
}

type Subscriber struct{}

func (s *Subscriber) EventCreated(ctx context.Context, id int) error {
	fmt.Printf("Received event #%d\n", id)
	return nil
}
