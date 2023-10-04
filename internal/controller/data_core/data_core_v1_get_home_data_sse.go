package data_core

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"home-network-watcher/api/data_core/v1"
	"home-network-watcher/manifest"
	"time"
)

func (c *ControllerV1) GetHomeDataSSE(ctx context.Context, _ *v1.GetHomeDataSSEReq) (_ *v1.GetHomeDataSSERes, err error) {
	request := g.RequestFromCtx(ctx)
	request.Response.Header().Set("Content-Type", "text/event-stream")
	request.Response.Header().Set("Cache-Control", "no-cache")
	request.Response.Header().Set("Connection", "keep-alive")
	request.Response.Header().Set("Access-Control-Allow-Origin", "*")
	request.Response.Header().Set("X-Accel-Buffering", "no")

	for {
		// 从缓存中获取数据
		homeData := gcache.MustGet(ctx, manifest.HaEntitiesCacheKey)
		resJson := gjson.New(&v1.GetHomeDataSSERes{HomeData: homeData}).MustToJsonString()

		// 发送数据
		request.Response.Writefln("data: " + resJson + "\n")
		request.Response.Flush()

		// 等待1秒或者上下文取消
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			glog.Info(ctx, "GetHomeDataSSE: ctx.Done()")
			request.ExitAll()
			return nil, ctx.Err()
		}
	}
}
