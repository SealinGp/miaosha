package ratelimiter

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"time"
)

var ErrLimitExceed = errors.New("rate limit exceed")

func NewTokenBucketLimiterWithBuildIn(bk *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bk.Allow() {
				return nil,ErrLimitExceed
			}
			return next(ctx,request)
		}
	}
}

//func a()  {
//	ratebucket := rate.NewLimiter(rate.Every(time.Second*1),100)
//}

func DynamicLimiter(interval,bust int) endpoint.Middleware {
	bucket := rate.NewLimiter(rate.Every(time.Second*time.Duration(interval)),bust)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bucket.Allow() {
				return nil,ErrLimitExceed
			}
			return next(ctx,request)
		}
	}
}