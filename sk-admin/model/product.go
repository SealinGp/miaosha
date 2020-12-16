package model

import (
	"github.com/gohouse/gorose/v2"
	"miaosha/pkg/mysql"
)

type Product struct {
	ProductModelId   int    `json:"product_id"`
	ProductModelName string `json:"product_name"`
	Total            int    `json:"total"`
	Status           int    `json:"status"`
}

type ProductModel struct {
}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) getTableName() string {
	return "product"
}

func (p *ProductModel) GetProductList() ([]gorose.Data, error) {
	return mysql.DB().Table(p.getTableName()).Get()
}

func (p *ProductModel) CreateProduct(product *Product) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"product_name": product.ProductModelName,
		"total":        product.Total,
		"status":       product.Status,
	}).Insert()
	return err
}
