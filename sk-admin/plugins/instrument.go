package plugins

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/gohouse/gorose/v2"
	"github.com/juju/ratelimit"
	"golang.org/x/time/rate"
	"miaosha/sk-admin/model"
	"miaosha/sk-admin/service"
	"time"
)

var ErrLimitExceed = errors.New("Rate limit exceed!")

func NewTokenBUcketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

type skAdminMetricMiddleware struct {
	service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type activityMetricsMiddleware struct {
	service.ActivityService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type productMetricsMiddleware struct {
	service.ProductService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (mw skAdminMetricMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}

func (mw productMetricsMiddleware) CreateProduct(product *model.Product) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.ProductService.CreateProduct(product)
}

func (mw productMetricsMiddleware) GetProductList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.ProductService.GetProductList()
}

func (mw activityMetricsMiddleware) GetActivityList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.ActivityService.GetActivityList()
}

func (mw activityMetricsMiddleware) CreateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.ActivityService.CreateActivity(activity)
}

func SkAdminMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAdminMetricMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

func ActivityMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityMetricsMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}
func ProductMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productMetricsMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}
