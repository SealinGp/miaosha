package model

import (
	"github.com/gohouse/gorose/v2"
	"miaosha/pkg/mysql"
)

type Product struct {
	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Total       int    `json:"total"`
	Status      int    `json:"status"`
}

type ProductMap struct {
}

func (p *Product) getTableName() string {
	return "product"
}

func (p *Product) GetProductList() ([]gorose.Data, error) {
	return mysql.DB().Table(p.getTableName()).Get()
}

func (p *Product) CreateProduct(product *Product) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"product_name": product.ProductName,
		"total":        product.Total,
		"status":       product.Status,
	}).Insert()
	return err
}
