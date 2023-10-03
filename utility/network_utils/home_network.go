package network_utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/hamster1963/360-router-data-retriever/rconfig"
	"github.com/hamster1963/360-router-data-retriever/rutils"
	"home-network-watcher/manifest"
	"time"
)

type uNetworkUtils struct {
}

var NetworkUtils = &uNetworkUtils{}

var (
	routerConfig = &rconfig.RouterConfig{
		RouterIP:       manifest.HomeNetworkRouterIP,
		RouterAddress:  manifest.HomeNetworkRouterAddress,
		RouterPassword: manifest.HomeNetworkRouterPassword,
	}
	routerMain rutils.SRouterController = rutils.NewRouter().InitRouter(routerConfig)
)

// GetHomeNetwork
//
//	@dc: 获取家庭路由器网速
//	@author: laixin   @date:2023/4/2 19:43:13
func (u *uNetworkUtils) GetHomeNetwork() (err error) {
	var (
		homeNetwork = g.Map{
			"time":        "",
			"rxSpeedKbps": 0,
			"txSpeedKbps": 0,
			"rxSpeedMbps": 0,
			"txSpeedMbps": 0,
		}
	)

	// 检测登陆状态
	if login, err := routerMain.CheckLogin(); err != nil || login == false {
		err := routerMain.GetRandomString()
		if err != nil {
			return err
		}
		err = routerMain.GenerateAesString()
		if err != nil {
			return err
		}
		err = routerMain.Login()
		if err != nil {
			return err
		}
	}
	routerSpeedInfo, err := routerMain.GetRouterSpeed()
	if err != nil {
		return err
	}

	jsonData := gjson.New(routerSpeedInfo)
	rxSpeed := jsonData.Get("data.down_speed") // 下载速度
	txSpeed := jsonData.Get("data.up_speed")   // 上传速度

	// 速度单位转换
	rxSpeedKbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(rxSpeed)/1024))
	txSpeedKbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(txSpeed)/1024))
	homeNetwork["rxSpeedKbps"] = rxSpeedKbps
	homeNetwork["txSpeedKbps"] = txSpeedKbps

	// 转换成MB
	rxSpeedMbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(rxSpeed)/1024/1024))
	txSpeedMbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(txSpeed)/1024/1024))
	homeNetwork["rxSpeedMbps"] = rxSpeedMbps
	homeNetwork["txSpeedMbps"] = txSpeedMbps

	homeNetwork["time"] = gtime.Now().String()
	deviceData, err := routerMain.GetDeviceList()
	if err != nil {
		return err
	}
	homeNetwork["deviceCount"] = len(gjson.New(deviceData).Get("client_node").Array())

	err = gcache.Set(context.Background(), manifest.HomeNetworkCacheKey, homeNetwork, 10*time.Second)
	if err != nil {
		return err
	}
	return nil
}
