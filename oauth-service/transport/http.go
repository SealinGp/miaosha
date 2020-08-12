package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	oauthServiceEndpoint "miaosha/oauth-service/endpoint"
	"miaosha/oauth-service/service"
	"net/http"
)

var (
	ErrBadRequest  = errors.New("invalid request parameter")
	ErrGrantTypeRequest = errors.New("invalid request grant type")
	ErrTokenRequest = errors.New("invalid request token")
	ErrInvalidClientRequest = errors.New("invalid client message")
)

func MakeHttpHandler(
	ctx context.Context,endpoints oauthServiceEndpoint.OAuth2Endpoints,
	tokenService service.TokenService,clientService service.ClientDetailsService,
	zipkinTracer *gozipkin.Tracer,
	logger log.Logger,
)  {
	r            := mux.NewRouter()
	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer,zipkin.Name("http-transport"))
	options      := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}
	r.Path("/metrics").Handler(promhttp.Handler())

	clientAuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeClientAuthorizationContext(clientService,logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Methods("POST").Path("/oauth/token").Handler(kithttp.NewServer(
		endpoints.TokenEndpoint,
		decodeTokenRequest,
		encodeJsonResponse,
		clientAuthorizationOptions...
	))

	r.Methods("POST").Path("/oauth/check_token").Handler(kithttp.NewServer(
		endpoints.CheckTokenEndpoint,
		decodeCheckTokenRequest,
		encodeJsonResponse,
		clientAuthorizationOptions...
	))
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...
	))
}

func makeClientAuthorizationContext(clientDetailsService service.ClientDetailsService,logger log.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		if clientId, clientSecret, ok := r.BasicAuth(); ok {
			clientDetails, err := clientDetailsService.GetClientDetailsByClientId(ctx,clientId,clientSecret)
			if err == nil {
				return context.WithValue(ctx,oauthServiceEndpoint.OAuth2ClientDetailsKey,clientDetails)
			}
		}
		return context.WithValue(ctx,oauthServiceEndpoint.OAuth2ErrorKey,ErrInvalidClientRequest)
	}
}

func decodeTokenRequest(ctx context.Context,r *http.Request) (interface{},error) {
	grantType := r.URL.Query().Get("grant_type")
	if grantType == "" {
		return nil,ErrGrantTypeRequest
	}
	return &oauthServiceEndpoint.TokenRequest{
		GrantType: grantType,
		Reader:    r,
	},nil
}

func decodeCheckTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	tokenValue := r.URL.Query().Get("token")
	if tokenValue == "" {
		return nil,ErrTokenRequest
	}
	return &oauthServiceEndpoint.CheckTokenRequest{
		Token:         tokenValue,
	},nil
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter)  {
	w.Header().Set("Content-Type","application/json;charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"err":err.Error(),
	})
}
func encodeJsonResponse(ctx context.Context,w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type","application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func decodeHealthCheckRequest(ctx context.Context,r *http.Request) (interface{},error) {
	return oauthServiceEndpoint.HealthRequest{},nil
}