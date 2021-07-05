package service

import (
	"github.com/gohouse/gorose/v2"
	"miaosha/sk-admin/model"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductList() ([]gorose.Data, error)
}

type ProductServiceMiddleware func(service ProductService) ProductService

type ProductServiceImpl struct {
}

func (productServiceImpl *ProductServiceImpl) CreateProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	return productEntity.CreateProduct(product)
}

func (productServiceImpl *ProductServiceImpl) GetProductList() ([]gorose.Data, error) {
	productEntity := model.NewProductModel()
	return productEntity.GetProductList()
}
