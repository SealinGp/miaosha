package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/unknwon/com"
	"log"
	"miaosha/pkg/config"
	"miaosha/sk-admin/model"
	"time"
)

type ActivityService interface {
	GetActivityList() ([]gorose.Data, error)
	CreateActivity(activity *model.Activity) error
}

type ActivityServiceMiddleware func(ActivityService) ActivityService

type ActivityServiceImpl struct {
}

func (activityServiceImpl *ActivityServiceImpl) GetActivityList() ([]gorose.Data, error) {
	activityEntity := model.NewActivityModel()
	activityList, err := activityEntity.GetActivityList()
	if err != nil {
		return nil, err
	}
	for _, v := range activityList {
		startTime, _ := com.StrTo(fmt.Sprint(v["start_time"])).Int64()
		v["start_time_str"] = time.Unix(startTime, 0).Format("2006-01-02 15:04:05")

		endTime, _ := com.StrTo(fmt.Sprint(v["end_time"])).Int64()
		v["end_time_str"] = time.Unix(endTime, 0).Format("2006-01-02 15:04:05")

		nowTime := time.Now().Unix()
		if nowTime > endTime {
			v["status_str"] = "已结束"
			continue
		}

		status, _ := com.StrTo(fmt.Sprint(v["status"])).Int()
		if status == model.ActivityStatusNormal {
			v["status_str"] = "正常"
		} else if status == model.ActivityStatusDisable {
			v["status_str"] = "已禁用"
		}
	}
	return activityList, nil
}

func (activityServiceImpl *ActivityServiceImpl) CreateActivity(activity *model.Activity) error {
	activityEntity := model.NewActivityModel()
	err := activityEntity.CreateActivity(activity)
	if err != nil {
		return err
	}

	return activityServiceImpl.syncToZk(activity)
}

func (activityServiceImpl *ActivityServiceImpl) syncToZk(activity *model.Activity) error {
	zkPath := config.Zk.SecProductKey
	secProductInfoList, err := activityServiceImpl.loadProductFromZk(zkPath)
	if err != nil {
		return err
	}
	var secProductInfo = &model.SecProductInfoConf{}
	secProductInfo.EndTime = activity.EndTime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.SoldMaxLimit = activity.Speed
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.BuyRate = activity.BuyRate
	secProductInfoList = append(secProductInfoList, secProductInfo)

	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		return err
	}
	conn := config.Zk.ZkConn

	var byteData = data
	var flags int32
	var acls = zk.WorldACL(zk.PermAll)

	exists, _, _ := conn.Exists(zkPath)
	if exists {
		_, err := conn.Create(zkPath, byteData, flags, acls)
		if err != nil {
			//todo log
		}
	} else {
		_, err := conn.Create(zkPath, byteData, flags, acls)
		if err != nil {
			//todo log
		}
	}
	return nil
}

func (activityServiceImpl *ActivityServiceImpl) loadProductFromZk(key string) ([]*model.SecProductInfoConf, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	v, s, err := config.Zk.ZkConn.Get(key)
	if err != nil {
		return nil, err
	}
	log.Printf("[I] get from zk success. data:%v", s)

	var secProductInfo []*model.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfo)
	return secProductInfo, err
}
