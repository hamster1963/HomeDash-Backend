package data_core

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"home-network-watcher/manifest"
	"time"

	"home-network-watcher/api/data_core/v1"
)

func (c *ControllerV1) GetNetworkDataSSE(ctx context.Context, _ *v1.GetNetworkDataSSEReq) (_ *v1.GetNetworkDataSSERes, err error) {
	request := g.RequestFromCtx(ctx)
	request.Response.Header().Set("Content-Type", "text/event-stream")
	request.Response.Header().Set("Cache-Control", "no-cache")
	request.Response.Header().Set("Connection", "keep-alive")
	request.Response.Header().Set("Access-Control-Allow-Origin", "*")
	request.Response.Header().Set("X-Accel-Buffering", "no")

	for {
		// 从缓存中获取数据
		nodeInfo := gcache.MustGet(ctx, manifest.ProxyNodeCacheKey)
		homeNetwork := gcache.MustGet(ctx, manifest.HomeNetworkCacheKey)
		proxyNetwork := gcache.MustGet(ctx, manifest.ProxyNetworkCacheKey)
		coffeeInfo := gcache.MustGet(ctx, manifest.ProxySubscribeCacheKey)
		serverInfo := gcache.MustGet(ctx, manifest.ServerDataCacheKey)

		resJson := gjson.New(&v1.GetNetworkDataSSERes{
			NodeInfo:     nodeInfo,
			HomeNetwork:  homeNetwork,
			ProxyNetwork: proxyNetwork,
			CoffeeInfo:   coffeeInfo,
			ServerInfo:   serverInfo,
		}).MustToJsonString()

		// 发送数据
		request.Response.Writefln("data: " + resJson + "\n")
		request.Response.Flush()

		// 等待1秒或者上下文取消
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			glog.Info(ctx, "GetNetworkDataSSE: ctx.Done()")
			request.ExitAll()
			return nil, ctx.Err()
		}
	}

}
