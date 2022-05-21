package eventpubsub

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, typ string, data interface{}) error
}
