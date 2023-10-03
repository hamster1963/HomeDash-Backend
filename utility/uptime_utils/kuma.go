package uptime_utils

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"home-network-watcher/manifest"
	"time"
)

type uKuma struct{}

var Kuma = &uKuma{}

type UptimeInfo struct {
	Status bool   `json:"status"`
	Ping   int    `json:"ping"`
	Uptime string `json:"uptime"`
}

// GetMonitorStatus
//
//	@dc: 获取监控状态
//	@author: hamster   @date:2023/9/20 14:43:18
func (u *uKuma) GetMonitorStatus(ctx context.Context) (serverCount int, errorServer int, err error) {
	response, err := g.Client().Get(ctx, manifest.UptimeKumaApiUrl)
	defer func(response *gclient.Response) {
		err := response.Close()
		if err != nil {
			return
		}
	}(response)
	if err != nil {
		return 0, 0, err
	}
	// 获取返还数据
	heartBeatJson := gjson.New(response.ReadAllString()).Get("heartbeatList").Map()
	serverCount = len(heartBeatJson)
	for _, value := range heartBeatJson {
		pingList := gconv.SliceAny(value)
		if gjson.New(pingList[len(pingList)-1]).Get("status").Int() != 1 {
			errorServer++
		}
	}
	return serverCount, errorServer, nil
}

// GetUptimeData
//
//	@dc: 获取核心服务监控数据
//	@author: hamster   @date:2023/9/21 11:21:13
func (u *uKuma) GetUptimeData(ctx context.Context) (err error) {
	serviceMap := g.Map{
		"xui":     2,
		"v2raya":  9,
		"proxy":   8,
		"nginx":   5,
		"home":    7,
		"netflix": 15}
	response, err := g.Client().Get(ctx, manifest.UptimeKumaApiUrl)
	if err != nil {
		return err
	}
	defer func(response *gclient.Response) {
		err := response.Close()
		if err != nil {
			return
		}
	}(response)
	// 获取返还数据，获取最新状态与 24 小时状态
	heartBeatJson := gjson.New(response.ReadAllString())
	for key, value := range serviceMap {
		serviceUptime := &UptimeInfo{}
		statusList := heartBeatJson.Get("heartbeatList" + "." + gconv.String(value)).Array()
		if len(statusList) == 0 {
			serviceUptime.Status = false
			serviceUptime.Ping = 0
			serviceUptime.Uptime = "0.00"
			serviceMap[key] = gconv.Map(serviceUptime)
			continue
		}
		// 获取最后一条数据
		lastHeartBeat := statusList[len(statusList)-1]
		if gconv.Int(gconv.Map(lastHeartBeat)["status"]) == 1 {
			serviceUptime.Status = true
		} else {
			serviceUptime.Status = false
		}
		serviceUptime.Ping = gconv.Int(gconv.Map(lastHeartBeat)["ping"])
		// 获取24小时可用率
		uptimeList := heartBeatJson.Get("uptimeList" + "." + gconv.String(value) + "_24").Float64()
		switch uptimeList {
		case 0:
			serviceUptime.Uptime = "0.00"
		case 1:
			serviceUptime.Uptime = "100.0"
		default:
			// 转换为2位小数百分比
			serviceUptime.Uptime = gconv.String(uptimeList * 100)
			if len(serviceUptime.Uptime) > 4 {
				serviceUptime.Uptime = serviceUptime.Uptime[:4]
			}
		}
		serviceMap[key] = gconv.Map(serviceUptime)
	}
	err = gcache.Set(ctx, manifest.UptimeCacheKey, serviceMap, 1*time.Hour)
	if err != nil {
		return err
	}
	return
}
