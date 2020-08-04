package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"miaosha/oauth-service/model"
)
var (
	ErrInvalidClientRequest = errors.New("invalid client message")
	ErrInvalidUserRequest = errors.New("invalid user message")
	ErrNotPermit = errors.New("not permit")
)
//
const (
	OAuth2DetailsKey = "OAuth2Details"
	OAuth2ClientDetailsKey = "OAuth2ClientDetails"
	OAuth2ErrorKey = "OAuth2Error"
)

type OAuth2Endpoints struct {
	TokenEndpoint endpoint.Endpoint
	CheckTokenEndpoint endpoint.Endpoint
	GRPCCheckTokenEndpoint endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func MakeClientAuthorizationMiddleware(logger log.Logger) endpoint.Middleware {

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil,err
			}
			if _, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails); !ok {
				return nil,ErrInvalidClientRequest
			}
			return next(ctx,request)
		}
	}
}

func MakeOAuth2AuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil,err
			}
			if _, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); !ok  {
				return nil,ErrInvalidUserRequest
			}
			return next(ctx,request)
		}
	}
}

func MakeAuthorityAuthorizationMiddleware(authority string,logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil,err
			}
			details,ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details)
			if !ok {
				return nil,ErrInvalidClientRequest
			}
			hasAuthority := false
			for _, value := range details.User.Authorities {
				if value == authority {
					hasAuthority = true
					break
				}
			}

			if !hasAuthority {
				return nil,ErrNotPermit
			}
			return next(ctx,request)
		}
	}
}

type TokenRequest struct {
	GrantType string

}