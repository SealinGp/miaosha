package main

import (
	"miaosha/pkg/config"
	"miaosha/pkg/mysql"
	"miaosha/sk-admin/setup"
)

func main() {
	mysql.InitMysql(
		config.MysqlConfig.Host,
		config.MysqlConfig.Port,
		config.MysqlConfig.User,
		config.MysqlConfig.Pwd,
		config.MysqlConfig.Db,
	)
	setup.InitZk()

}
