package eventpubsub

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

var (
	ErrInvalidType = errors.New("invalid event type")
	ErrInvalidData = errors.New("invalid event data")
)

type Event interface {
	Type() string
	Data() interface{}
}

type Handler interface {
	Handle(ctx context.Context, event Event) error
}

type HandlerSet struct {
	set map[string]Handler
}

func NewHandlerSet() *HandlerSet {
	return &HandlerSet{set: make(map[string]Handler)}
}

func (hs *HandlerSet) Add(typ string, handler Handler) {
	hs.set[typ] = handler
}

func (hs *HandlerSet) Handle(ctx context.Context, event Event) error {
	h, ok := hs.set[event.Type()]
	if !ok {
		return ErrInvalidType
	}
	return h.Handle(ctx, event)
}

// DecodeInputFunc extracts a user-domain input object from an event.
// It's designed to be used in event handlers, for subscriber-side endpoints.
type DecodeInputFunc func(context.Context, Event) (input interface{}, err error)

// Subscriber wraps an endpoint and implements Handler.
type Subscriber struct {
	e   endpoint.Endpoint
	dec DecodeInputFunc
}

// NewSubscriber constructs a new subscriber, which implements Handler
// and wraps the provided endpoint.
func NewSubscriber(e endpoint.Endpoint, dec DecodeInputFunc) *Subscriber {
	return &Subscriber{
		e:   e,
		dec: dec,
	}
}

// Handle implements Handler.
func (s *Subscriber) Handle(ctx context.Context, event Event) error {
	input, err := s.dec(ctx, event)
	if err != nil {
		return err
	}

	output, _ := s.e(ctx, input) // err is always nil
	if f, ok := output.(endpoint.Failer); ok {
		return f.Failed()
	}
	return nil
}
