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

func (c *ControllerV1) GetDockerMonitorSSE(ctx context.Context, _ *v1.GetDockerMonitorSSEReq) (_ *v1.GetDockerMonitorSSERes, err error) {
	request := g.RequestFromCtx(ctx)
	request.Response.Header().Set("Content-Type", "text/event-stream")
	request.Response.Header().Set("Cache-Control", "no-cache")
	request.Response.Header().Set("Connection", "keep-alive")
	request.Response.Header().Set("Access-Control-Allow-Origin", "*")
	request.Response.Header().Set("X-Accel-Buffering", "no")

	for {
		// 从缓存中获取数据
		dockerData := gcache.MustGet(ctx, manifest.DockerMonitorCacheKey)
		resJson := gjson.New(&v1.GetDockerMonitorSSERes{DockerData: dockerData}).MustToJsonString()

		// 发送数据
		request.Response.Writefln("data: " + resJson + "\n")
		request.Response.Flush()

		// 等待10秒或者上下文取消
		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
			glog.Info(ctx, "GetDockerMonitorSSE: ctx.Done()")
			request.ExitAll()
			return nil, ctx.Err()
		}
	}
}
