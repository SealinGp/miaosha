package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

var engin *gorose.Engin
var err error
func InitMysql(host,port,user,pwd,db string)  {
	DbConfig := gorose.Config{
		Driver:"mysql",
		Dsn:fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",user,pwd,host,port,db),
		Prefix:"",
		SetMaxOpenConns:300,
		SetMaxIdleConns:10,
	}
	engin, err = gorose.Open(DbConfig)
	if err != nil {
		fmt.Printf(err)
		return
	}
}
func DB() gorose.IOrm {
	return engin.NewOrm()
}