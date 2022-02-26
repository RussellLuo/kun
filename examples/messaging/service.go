package messaging

import (
	"context"
	"fmt"
)

//go:generate kungen ./service.go Service

type Service interface {
	//kun:op GET /messages/{messageID}
	//kun:op GET /users/{userID}/messages/{messageID}
	GetMessage(ctx context.Context, userID string, messageID string) (text string, err error)
}

type Messaging struct{}

func (m *Messaging) GetMessage(ctx context.Context, userID string, messageID string) (string, error) {
	return fmt.Sprintf("user[%s]: message[%s]", userID, messageID), nil
}
