package helloworldgrpc

import (
	"context"
)

//go:generate kokgen ./service.go Service

// Service is used for saying hello.
type Service interface {
	// SayHello says hello to the given name.
	//kok:grpc
	SayHello(ctx context.Context, name string) (message string, err error)
}

type Greeter struct{}

func (g *Greeter) SayHello(ctx context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
