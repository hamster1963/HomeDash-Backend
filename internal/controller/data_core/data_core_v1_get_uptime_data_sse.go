package data_core

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"home-network-watcher/api/data_core/v1"
	"home-network-watcher/manifest"
	"time"
)

func (c *ControllerV1) GetUptimeDataSSE(ctx context.Context, _ *v1.GetUptimeDataSSEReq) (res *v1.GetUptimeDataSSERes, err error) {
	request := g.RequestFromCtx(ctx)
	request.Response.Header().Set("Content-Type", "text/event-stream")
	request.Response.Header().Set("Cache-Control", "no-cache")
	request.Response.Header().Set("Connection", "keep-alive")
	request.Response.Header().Set("Access-Control-Allow-Origin", "*")
	request.Response.Header().Set("X-Accel-Buffering", "no")

	for {
		// 从缓存中获取数据
		uptimeData, err := gcache.Get(ctx, manifest.UptimeCacheKey)
		if err != nil {
			return nil, err
		}
		res = &v1.GetUptimeDataSSERes{UptimeData: uptimeData}
		// 发送数据
		request.Response.Writefln("data: " + gjson.New(res).MustToJsonString() + "\n")
		request.Response.Flush()

		// 等待10秒或者上下文取消
		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
