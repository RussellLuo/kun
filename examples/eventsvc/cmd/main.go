package main

import (
	"context"

	"github.com/RussellLuo/kun/examples/eventsvc"
	"github.com/RussellLuo/kun/pkg/eventcodec"
	"github.com/RussellLuo/kun/pkg/eventpubsub"
)

type event struct {
	typ  string
	data []byte
}

func (e *event) Type() string {
	return e.typ
}

func (e *event) Data() interface{} {
	return e.data
}

type publisher struct {
	sub eventpubsub.Handler
}

func (p *publisher) Publish(ctx context.Context, typ string, data interface{}) error {
	e := &event{
		typ:  typ,
		data: data.([]byte),
	}

	// For demonstration, we deliver the event e to the subscriber sub directly.
	return p.sub.Handle(ctx, e)
}

func main() {
	codecs := eventcodec.NewDefaultCodecs(nil)
	sub := eventsvc.NewEventHandler(&eventsvc.Subscriber{}, codecs)

	pub := eventsvc.NewEventPublisher(codecs, &publisher{sub: sub})
	if err := pub.EventCreated(context.Background(), 1); err != nil {
		panic(err)
	}
}
