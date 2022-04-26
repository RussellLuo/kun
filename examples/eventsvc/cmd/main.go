package main

import (
	"context"

	"github.com/RussellLuo/kun/examples/eventsvc"
	"github.com/RussellLuo/kun/pkg/eventcodec"
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

func main() {
	svc := &eventsvc.Subscriber{}
	handler := eventsvc.NewEventHandler(svc, eventcodec.NewDefaultCodecs(nil))

	e := &event{
		typ:  "created",
		data: []byte(`{"id": 1}`),
	}
	if err := handler.Handle(context.Background(), e); err != nil {
		panic(err)
	}
}
