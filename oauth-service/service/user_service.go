package service

import (
	"context"
	"errors"
	"miaosha/oauth-service/model"
	"miaosha/pb"
	"miaosha/pkg/client"
)
var (
	InvalidAuthentication = errors.New("invalid auth")
	InvalidUserInfo = errors.New("invalid user info")
)
type UserDetailService interface {
	GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails,error)
}

type RemoteUserService struct {
	userClient client.UserClient
}

func (service *RemoteUserService)GetUserDetailByUsername(ctx context.Context,username, password string) (*model.UserDetails,error) {
	response,err := service.userClient.CheckUser(ctx,nil,&pb.UserRequest{
		Username:username,
		Password:password,
	})
	if err != nil {
		return nil,err
	}
	if response.UserId == 0 {
		return nil,InvalidUserInfo
	}

	return &model.UserDetails{
		UserId:response.UserId,
		Username:username,
		Password:password,
	},nil
}

func NewRemoteUserService() *RemoteUserService {
	userClient, _ := client.NewUserClient("user",nil,nil)
	return &RemoteUserService{
		userClient:userClient,
	}
}


type ServiceMiddleware func(Service) Service