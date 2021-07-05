package main

import "miaosha/pkg/config"

func remoteCfgInit() {
	if err := config.SubParse("mysql", &config.MysqlConfig); err != nil {
		config.Logger.Log("Fail to parse mysql", err)
	}
}
