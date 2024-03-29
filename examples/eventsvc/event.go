// Code generated by kun; DO NOT EDIT.
// github.com/RussellLuo/kun

package eventsvc

import (
	"context"

	"github.com/RussellLuo/kun/pkg/eventcodec"
	"github.com/RussellLuo/kun/pkg/eventpubsub"
)

func NewEventHandler(svc Service, codecs eventcodec.Codecs) eventpubsub.Handler {
	var codec eventcodec.Codec
	handlerSet := eventpubsub.NewHandlerSet()

	codec = codecs.EncodeDecoder("EventCreated")
	handlerSet.Add("created", eventpubsub.NewSubscriber(
		MakeEndpointOfEventCreated(svc),
		decodeEventCreatedInput(codec),
	))

	return handlerSet
}

func decodeEventCreatedInput(codec eventcodec.Codec) eventpubsub.DecodeInputFunc {
	return func(_ context.Context, event eventpubsub.Event) (interface{}, error) {
		var input EventCreatedRequest

		if err := codec.Decode(event.Data(), &input); err != nil {
			return nil, err
		}

		return &input, nil
	}
}

// EventPublisher implements Service on the publisher side.
//
// EventPublisher should only be used in limited scenarios where only one subscriber
// is involved and the publisher depends on the interface provided by the subscriber.
//
// In typical use cases of the publish-subscribe pattern - many subscribers are
// involved and the publisher knows nothing about the subscribers - you should
// just send the event in the way it should be.
type EventPublisher struct {
	publisher eventpubsub.Publisher
	codecs    eventcodec.Codecs
}

func NewEventPublisher(publisher eventpubsub.Publisher, codecs eventcodec.Codecs) *EventPublisher {
	return &EventPublisher{
		publisher: publisher,
		codecs:    codecs,
	}
}

func (p *EventPublisher) EventCreated(ctx context.Context, id int) (err error) {
	codec := p.codecs.EncodeDecoder("EventCreated")

	_data, err := codec.Encode(&EventCreatedRequest{
		Id: id,
	})
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, "created", _data)
}
