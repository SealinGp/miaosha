package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"miaosha/oauth-service/model"
	"miaosha/oauth-service/service"
	"net/http"
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

//客户端认证信息
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

//OAuth2认证信息
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

//权限认证
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
	Reader *http.Request
}
type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
	Error string `json:"error"`
}
//获取token端点
func MakeTokenEndpoint(svc service.TokenGranter,clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req       := request.(*TokenRequest)
		token,err := svc.Grant(ctx,req.GrantType,ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails),req.Reader)
		errStr    := ""
		if err != nil {
			errStr = err.Error()
		}
		return TokenResponse{
			AccessToken: token,
			Error:       errStr,
		},nil
	}
}

type CheckTokenRequest struct {
	Token string
	ClientDetails model.ClientDetails
}
type CheckTokenResponse struct {
	OAuth2Details *model.OAuth2Details `json:"o_auth_2_details"`
	Error string `json:"error"`
}
//检查token端点
func MakeCheckTokenEndpoint(svc service.TokenService)  endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CheckTokenRequest)
		tokenDetails,err := svc.GetOAuth2DetailsByAccessToken(req.Token)
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		return CheckTokenResponse{
			OAuth2Details: tokenDetails,
			Error:         errStr,
		},nil
	}
}

type SimpleRequest struct {}
type SimpleResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}
type AdminRequest struct {}
type AdminResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

type HealthRequest struct {}
type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthResponse{Status:svc.HealthCheck()},nil
	}
}