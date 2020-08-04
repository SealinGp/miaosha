package transport

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
)

var (
	ErrBadRequest  = errors.New("invalid request parameter")
	ErrGrantTypeRequest = errors.New("invalid request grant type")
	ErrTokenRequest = errors.New("invalid request token")
	ErrInvalidClientRequest = errors.New("invalid client message")
)

func MakeHttpHandler(ctx context.Context,)  {
	mux.NewRouter()
}