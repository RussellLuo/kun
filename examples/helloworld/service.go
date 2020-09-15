package helloworld

import (
	"context"
)

//go:generate kokgen ./service.go Service

type Service interface {
	// @kok2(op): POST /messages
	SayHello(ctx context.Context, name string) (message string, err error)
}

type Greeter struct{}

func (g *Greeter) SayHello(ctx context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
