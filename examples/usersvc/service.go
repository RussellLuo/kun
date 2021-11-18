package usersvc

import (
	"context"
	"net"
)

//go:generate kokgen ./service.go Service

type User struct {
	Name string
	Age  int
	IP   net.IP `kok:"in=header name=X-Forwarded-For, in=request name=RemoteAddr"`
}

type Service interface {
	//kok:op POST /users
	//kok:param user
	//kok:success body=result
	CreateUser(ctx context.Context, user User) (result User, err error)
}

type UserService struct{}

func (u *UserService) CreateUser(ctx context.Context, user User) (User, error) {
	return user, nil
}
