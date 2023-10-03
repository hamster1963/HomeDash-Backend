package network_utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"home-network-watcher/manifest"
)

type uProxyProvider struct{}

var ProxyProvider = &uProxyProvider{}

// GetSubscribeInfo
//
//	@dc: 获取订阅信息
//	@author: laixin   @date:2023/6/6 04:15:01
func (u *uProxyProvider) GetSubscribeInfo() (err error) {
	err, authData := getAuthData()
	if err != nil || authData == "" {
		return
	}
	proxyClient := gclient.New()
	proxyClient.SetHeaderMap(map[string]string{
		"Authorization": authData,
	})
	response, err := proxyClient.Get(context.TODO(), manifest.ProxyProviderBaseUrl, nil)
	if err != nil {
		return
	}
	infoData := gconv.Map(gconv.Map(response.ReadAllString())["data"])
	if infoData["d"] == nil || infoData["transfer_enable"] == nil {
		return
	}
	usedBound := gconv.Float64(infoData["d"]) / 1024 / 1024 / 1010
	planBound := gconv.Float64(infoData["transfer_enable"]) / 1024 / 1024 / 1024
	remainBound := planBound - usedBound
	// 保留两位小数
	usedBoundStr := fmt.Sprintf("%.2f", usedBound)
	planBoundStr := fmt.Sprintf("%.2f", planBound)
	remainBoundStr := fmt.Sprintf("%.2f", remainBound)
	proxyCache := g.Map{
		"usedBound":   usedBoundStr,
		"planBound":   planBoundStr,
		"remainBound": remainBoundStr,
	}
	err = gcache.Set(context.TODO(), manifest.ProxySubscribeCacheKey, proxyCache, 0)
	if err != nil {
		return err
	}
	return
}

// getAuthData
//
//	@dc: 获取代理提供商AuthData
func getAuthData() (err error, authData string) {
	url := manifest.ProxyProviderLoginUrl
	response, err := gclient.New().Post(context.TODO(), url, manifest.ProxyProviderAuthData)
	if err != nil {
		return
	}
	return nil, gjson.New(response.ReadAllString()).Get("data.auth_data").String()
}
