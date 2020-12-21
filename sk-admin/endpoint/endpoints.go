package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/gohouse/gorose/v2"
	"log"
	"miaosha/sk-admin/model"
	"miaosha/sk-admin/service"
)

type SkAdminEndpoints struct {
	GetActivityEndpoint    endpoint.Endpoint
	CreateActivityEndpoint endpoint.Endpoint
	CreateProductEndpoint  endpoint.Endpoint
	GetProductEndpoint     endpoint.Endpoint
	HealthCheckEndpoint    endpoint.Endpoint
}

func (ue SkAdminEndpoints) HealthCheck() bool {
	return false
}

var (
	ErrInvalidRequestType = errors.New("invalid username,password")
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
}

type GetResponse struct {
	Result []gorose.Data `json:"result"`
	Error  error         `json:"error"`
}

type CreateResponse struct {
	Error error `json:"error"`
}

//  make endpoint
func MakeGetActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		log.Printf("GetActivityList")
		activityList, calError := svc.GetActivityList()
		if calError != nil {
			return GetResponse{Result: nil, Error: calError}, nil
		}
		return GetResponse{Result: activityList, Error: calError}, nil
	}
}

func MakeCreateActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*model.Activity)
		cErr := svc.CreateActivity(req)
		return CreateResponse{Error: cErr}, nil
	}
}

//  make endpoint
func MakeGetProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getProductList, calError := svc.GetProductList()
		if calError != nil {
			return GetResponse{Result: nil, Error: calError}, nil
		}
		return GetResponse{Result: getProductList, Error: calError}, nil
	}
}

func MakeCreateProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*model.Product)
		cErr := svc.CreateProduct(req)
		return CreateResponse{cErr}, nil
	}
}

type HealthRequest struct {
}

type GetListRequest struct{}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{status}, nil
	}
}
