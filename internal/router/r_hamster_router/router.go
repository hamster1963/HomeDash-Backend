package r_hamster_router

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"home-network-watcher/internal/controller/data_core"
)

func BindController(group *ghttp.RouterGroup) {
	group.Group("/", func(group *ghttp.RouterGroup) {
		BindDataCore(group)
	})
}

// BindDataCore 注册核心数据路由
func BindDataCore(group *ghttp.RouterGroup) {
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(data_core.NewV1())
	})
}
