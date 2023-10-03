package network_utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"home-network-watcher/manifest"
	"time"
)

type uProxyNetwork struct{}

var ProxyNetwork = &uProxyNetwork{}

// GetXuiUserList
//
//	@dc: 获取xui用户列表
//	@author: hamster   @date:2023/9/19 09:29:33
func (u *uProxyNetwork) GetXuiUserList() (err error) {
	var (
		url     = manifest.XuiBaseUrl + "/xui/inbound/list"
		session map[string]string
	)
	// 尝试获取缓存中的session
	sessionData, err := gcache.Get(context.Background(), manifest.ProxySessionCacheKey)
	if err != nil {
		return err
	}

	if sessionData.IsNil() {
		// 重新获取session
		session, err = u.GetXuiSession()
		if err != nil {
			return err
		}
	} else {
		session = sessionData.MapStrStr()
	}
	if err != nil {
		return err
	}
	post, err := g.Client().SetCookieMap(session).Post(context.Background(), url)
	defer func(post *gclient.Response) {
		err := post.Close()
		if err != nil {
			glog.Warning(context.Background(), err)
		}
	}(post)
	if err != nil {
		return err
	}
	if post.StatusCode != 200 {
		glog.Warning(context.Background(), "获取xui用户列表失败")
		return err
	}
	jsonData := gjson.New(post.ReadAllString())
	userList := jsonData.Get("obj").Array()
	// 删除敏感信息,计算总上传下载流量
	userCacheList := make([]g.Map, 0)
	upTotal := 0.00
	downTotal := 0.00
	for _, value := range userList {
		userJson := gjson.New(value)
		userCacheJson := g.Map{
			"id":       userJson.Get("id").Int(),
			"enable":   userJson.Get("enable").Bool(),
			"protocol": userJson.Get("protocol").String(),
			"remark":   userJson.Get("remark").String(),
			// 转换为 GB
			"up":   userJson.Get("up").Float64() / 1024 / 1024 / 1024,
			"down": userJson.Get("down").Float64() / 1024 / 1024 / 1024,
		}
		userCacheList = append(userCacheList, userCacheJson)
		upTotal += userJson.Get("up").Float64()
		downTotal += userJson.Get("down").Float64()
	}
	// 按照下载流量排序
	for i := 0; i < len(userCacheList); i++ {
		for j := i + 1; j < len(userCacheList); j++ {
			if userCacheList[i]["down"].(float64) < userCacheList[j]["down"].(float64) {
				userCacheList[i], userCacheList[j] = userCacheList[j], userCacheList[i]
			}
		}
	}
	cacheXuiData := g.Map{
		"user_list":  userCacheList,
		"user_count": len(userCacheList),
		"up_total":   upTotal / 1024 / 1024 / 1024,
		"down_total": downTotal / 1024 / 1024 / 1024,
	}
	err = gcache.Set(context.Background(), manifest.XUIUserListCacheKey, cacheXuiData, 1*time.Minute)
	return
}

// GetProxyNetwork
//
//	@dc: 获取代理服务器的网速
//	@author: laixin   @date:2023/4/2 20:06:21
func (u *uProxyNetwork) GetProxyNetwork() (err error) {
	var (
		proxyNetwork = g.Map{
			"time":        "",
			"rxSpeedKbps": 0,
			"txSpeedKbps": 0,
			"rxSpeedMbps": 0,
			"txSpeedMbps": 0,
		}
		session map[string]string
		url     = manifest.XuiBaseUrl + "/server/status"
	)

	// 尝试获取缓存中的session
	sessionData, err := gcache.Get(context.Background(), manifest.ProxySessionCacheKey)
	if err != nil {
		return err
	}

	if sessionData.IsNil() {
		// 重新获取session
		session, err = u.GetXuiSession()
		if err != nil {
			return err
		}
	} else {
		session = sessionData.MapStrStr()
	}

	// 通过xui进行网速的获取
	post, err := g.Client().SetCookieMap(session).Post(context.Background(), url)
	defer func(post *gclient.Response) {
		err := post.Close()
		if err != nil {
			glog.Warning(context.Background(), err)
		}
	}(post)
	if err != nil {
		return err
	}
	if post.StatusCode != 200 {
		glog.Warning(context.Background(), "获取网速失败")
		return err
	}
	jsonData := gjson.New(post.ReadAllString())
	rxSpeed := jsonData.Get("obj.netIO.down") // 下载速度
	txSpeed := jsonData.Get("obj.netIO.up")   // 上传速度

	// 速度单位转换
	rxSpeedKbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(rxSpeed)/1024))
	txSpeedKbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(txSpeed)/1024))
	proxyNetwork["rxSpeedKbps"] = rxSpeedKbps
	proxyNetwork["txSpeedKbps"] = txSpeedKbps

	// 转换成MB
	rxSpeedMbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(rxSpeed)/1024/1024))
	txSpeedMbps := gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(txSpeed)/1024/1024))
	proxyNetwork["rxSpeedMbps"] = rxSpeedMbps
	proxyNetwork["txSpeedMbps"] = txSpeedMbps

	proxyNetwork["time"] = gtime.Now().String()
	err = gcache.Set(context.Background(), manifest.ProxyNetworkCacheKey, proxyNetwork, 10*time.Second)
	if err != nil {
		return err
	}

	return err
}

// GetXuiSession
//
//	@dc: 获取Xui登陆session
//	@author: laixin   @date:2023/4/2 20:06:21
func (u *uProxyNetwork) GetXuiSession() (sessionMap map[string]string, err error) {
	var (
		url = manifest.XuiBaseUrl + "/login"
	)
	post, err := g.Client().Post(context.Background(), url, manifest.XuiLoginDataMap)
	defer func(post *gclient.Response) {
		err := post.Close()
		if err != nil {
			glog.Warning(context.Background(), err)
		}
	}(post)
	if err != nil {
		return nil, err
	}
	if post.StatusCode != 200 {
		return nil, fmt.Errorf("登录失败")
	}
	if post.Header.Get("Set-Cookie") == "" {
		return nil, fmt.Errorf("获取Cookie失败")
	}
	// 将session存入缓存
	err = gcache.Set(context.Background(), manifest.ProxySessionCacheKey, post.GetCookieMap(), 15*time.Minute)
	if err != nil {
		return nil, err
	}
	return post.GetCookieMap(), nil
}
