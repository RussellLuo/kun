package simple

import (
	"context"
)

type Counter interface {
	One(ctx context.Context, i int) (err error)
	Two(ctx context.Context, b bool, s ...string) (err error)
}
