package service

import (
	"context"
	"errors"
	"miaosha/oauth-service/model"
)

var (
	ErrClientMessage = errors.New("invalid client")
)
type ClientDetailsService interface {
	GetClientDetailsByClientId(ctx context.Context, clientId,clientSecret string)(*model.ClientDetails,error)
}

type MysqlClientDetailsService struct {}

func NewMysqlClientDetailsService () ClientDetailsService {
	return &MysqlClientDetailsService{}
}
func (service *MysqlClientDetailsService)GetClientDetailsByClientId(ctx context.Context, clientId,clientSecret string)(*model.ClientDetails,error) {
	clientDetailsModel := model.NewClientDetailsModel()
	details,err := clientDetailsModel.GetClientDetailsByClientId(clientId)
	if err != nil {
		return nil,err
	}
	if details.ClientSecret != clientSecret {
		return nil,ErrClientMessage
	}
	return details,nil
}