package server_utils

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"home-network-watcher/manifest"
	"math"
)

type uNezha struct{}

var Nezha = &uNezha{}

// GetAllServerInfo
//
//	@dc: 获取所有服务器信息
//	@author: hamster   @date:2023/9/18 14:37:42
func (u *uNezha) GetAllServerInfo(ctx context.Context) (err error) {
	// 1. 获取所有服务器信息
	//TODO (hamster) 2023/10/3: 优化获取服务器列表序号的方式
	serverList := []int{1, 2, 3, 4}
	nezhaClient := g.Client().SetHeader("Authorization", manifest.NezhaApiKey)
	serverDataList := make([]g.Map, 0)
	for _, value := range serverList {
		nezhaApi := manifest.NezhaApiUrl + gconv.String(value)
		response, err := nezhaClient.Get(ctx, nezhaApi)
		if err != nil {
			return err
		}
		// 整理数据
		cacheJson := g.Map{}
		resJsonData := gjson.New(response.ReadAllString()).Get("result.0").Map()
		serverJsonData := gjson.New(resJsonData)
		cacheJson["id"] = serverJsonData.Get("id").Int()
		cacheJson["name"] = serverJsonData.Get("name").String()
		// CPU 占用率保留两位小数
		cpuValue := serverJsonData.Get("status.CPU").Float64()
		roundedValue := math.Round(cpuValue*100) / 100
		cacheJson["cpu"] = roundedValue
		// 内存占用率保留两位小数
		memoryValue := serverJsonData.Get("host.MemTotal").Float64()
		memoryUsedValue := serverJsonData.Get("status.MemUsed").Float64()
		memoryUsedPercent := math.Round(memoryUsedValue/memoryValue*10000) / 100
		cacheJson["memory"] = memoryUsedPercent
		// 磁盘占用率保留两位小数
		diskValue := serverJsonData.Get("host.DiskTotal").Float64()
		diskUsedValue := serverJsonData.Get("status.DiskUsed").Float64()
		diskUsedPercent := math.Round(diskUsedValue/diskValue*10000) / 100
		cacheJson["disk"] = diskUsedPercent
		// 在线时长,转换为天
		uptimeValue := serverJsonData.Get("status.Uptime").Float64()
		uptimeDay := math.Round(uptimeValue/86400*100) / 100
		cacheJson["uptime"] = uptimeDay
		serverDataList = append(serverDataList, cacheJson)
		_ = response.Close()
		// 在线状态，通过last_active判断
		lastActiveValue := serverJsonData.Get("last_active").Int64()
		if gtime.Now().Unix()-lastActiveValue > 300 {
			cacheJson["status"] = "offline"
		} else {
			cacheJson["status"] = "online"
		}
	}
	err = gcache.Set(ctx, manifest.ServerDataCacheKey, serverDataList, 0)
	if err != nil {
		return err
	}
	return
}
