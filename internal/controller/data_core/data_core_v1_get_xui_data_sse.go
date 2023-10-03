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

func (c *ControllerV1) GetXuiDataSSE(ctx context.Context, _ *v1.GetXuiDataSSEReq) (res *v1.GetXuiDataSSERes, err error) {
	request := g.RequestFromCtx(ctx)
	request.Response.Header().Set("Content-Type", "text/event-stream")
	request.Response.Header().Set("Cache-Control", "no-cache")
	request.Response.Header().Set("Connection", "keep-alive")
	request.Response.Header().Set("Access-Control-Allow-Origin", "*")
	request.Response.Header().Set("X-Accel-Buffering", "no")

	for {
		// 从缓存中获取数据
		xuiData, err := gcache.Get(context.Background(), manifest.XUIUserListCacheKey)
		if err != nil {
			return nil, err
		}
		res = &v1.GetXuiDataSSERes{XuiData: xuiData}
		// 发送数据
		request.Response.Writefln("data: " + gjson.New(res).MustToJsonString() + "\n")
		request.Response.Flush()
		time.Sleep(5 * time.Second)
	}

}
