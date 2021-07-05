package main

import (
	"miaosha/pkg/bootstrap"
	"miaosha/pkg/config"
	"miaosha/pkg/mysql"
	"miaosha/sk-admin/setup"
)

func main() {
	remoteCfgInit()
	mysql.InitMysql(
		config.MysqlConfig.Host,
		config.MysqlConfig.Port,
		config.MysqlConfig.User,
		config.MysqlConfig.Pwd,
		config.MysqlConfig.Db,
	)
	setup.InitZk()
	setup.InitServer(bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port)
}
