package data_core

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"home-network-watcher/manifest"
	"time"

	"home-network-watcher/api/data_core/v1"
)

func (c *ControllerV1) GetNetworkDataSSE(ctx context.Context, _ *v1.GetNetworkDataSSEReq) (res *v1.GetNetworkDataSSERes, err error) {
	request := g.RequestFromCtx(ctx)
	request.Response.Header().Set("Content-Type", "text/event-stream")
	request.Response.Header().Set("Cache-Control", "no-cache")
	request.Response.Header().Set("Connection", "keep-alive")
	request.Response.Header().Set("Access-Control-Allow-Origin", "*")
	request.Response.Header().Set("X-Accel-Buffering", "no")

	for {
		// 从缓存中获取数据
		nodeInfo, err := gcache.Get(ctx, manifest.ProxyNodeCacheKey)
		if err != nil {
			return nil, err
		}
		homeNetwork, err := gcache.Get(ctx, manifest.HomeNetworkCacheKey)
		if err != nil {
			return nil, err
		}
		proxyNetwork, err := gcache.Get(ctx, manifest.ProxyNetworkCacheKey)
		if err != nil {
			return nil, err
		}
		coffeeInfo, err := gcache.Get(ctx, manifest.ProxySubscribeCacheKey)
		if err != nil {
			return nil, err
		}
		serverInfo, err := gcache.Get(ctx, manifest.ServerDataCacheKey)
		if err != nil {
			return nil, err
		}
		res = &v1.GetNetworkDataSSERes{
			NodeInfo:     nodeInfo,
			HomeNetwork:  homeNetwork,
			ProxyNetwork: proxyNetwork,
			CoffeeInfo:   coffeeInfo,
			ServerInfo:   serverInfo,
		}
		// 发送数据
		request.Response.Writefln("data: " + gjson.New(res).MustToJsonString() + "\n")
		request.Response.Flush()

		// 等待1秒或者上下文取消
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

}
