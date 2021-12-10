package usersvc

import (
	"context"
	"net"
)

//go:generate kungen ./service.go Service

type User struct {
	Name string
	Age  int
	IP   net.IP `kun:"in=header name=X-Forwarded-For, in=request name=RemoteAddr"`
}

type Service interface {
	//kun:op POST /users
	//kun:param user
	//kun:success body=result
	CreateUser(ctx context.Context, user User) (result User, err error)
}

type UserService struct{}

func (u *UserService) CreateUser(ctx context.Context, user User) (User, error) {
	return user, nil
}
