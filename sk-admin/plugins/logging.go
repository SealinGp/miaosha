package plugins

import (
	"github.com/go-kit/kit/log"
	"github.com/gohouse/gorose/v2"
	"miaosha/sk-admin/model"
	"miaosha/sk-admin/service"
	"time"
)

type skAdminLoggingMiddleware struct {
	service.Service
	logger log.Logger
}

type activityLoggingMiddleware struct {
	service.ActivityService
	logger log.Logger
}

type productLoggingMiddleware struct {
	service.ProductService
	logger log.Logger
}

func SkAdminLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAdminLoggingMiddleware{next, logger}
	}
}

func ActivityLoggingMiddleware(logger log.Logger) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityLoggingMiddleware{next, logger}
	}
}

func ProductLoggingMiddleware(logger log.Logger) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productLoggingMiddleware{next, logger}
	}
}

func (mw productLoggingMiddleware) CreateProduct(product *model.Product) error {
	defer func(begin time.Time) {
		_ = mw.logger.Log("func", "check", "product", product, time.Since(begin))
	}(time.Now())
	return mw.ProductService.CreateProduct(product)
}

func (mw productLoggingMiddleware) GetProductList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log("func", "check", "took", time.Since(begin))
	}(time.Now())
	return mw.ProductService.GetProductList()
}

func (mw activityLoggingMiddleware) CreateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		_ = mw.logger.Log("func", "Check", "result", "result", "took", time.Since(begin))
	}(time.Now())
	return mw.ActivityService.CreateActivity(activity)
}

func (mw skAdminLoggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		_ = mw.logger.Log("func", "HealthCheck", "result", result, "took", time.Since(begin))
	}(time.Now())
	result = mw.Service.HealthCheck()
	return result
}
