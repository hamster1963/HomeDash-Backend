package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gres"
	"home-network-watcher/internal/boot"
	"home-network-watcher/internal/global/g_functions"
	"home-network-watcher/internal/global/g_index"
	"home-network-watcher/internal/global/g_middleware"
	"home-network-watcher/internal/router/r_hamster_router"
	binInfo "home-network-watcher/utility/bin_utils"
	"runtime"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			gres.Dump()
			s := g.Server()

			// 性能分析
			runtime.SetMutexProfileFraction(1) // (非必需)开启对锁调用的跟踪
			runtime.SetBlockProfileRate(1)     // (非必需)开启对阻塞操作的跟踪
			s.EnablePProf()

			// 统一日志服务
			g_functions.SetDefaultHandler()
			// 服务状态码处理
			g_middleware.SMiddlewares.ErrorsStatus(s)

			// 全局中间件
			s.BindMiddlewareDefault(
				g_middleware.SMiddlewares.MiddlewareCORS,
				g_middleware.SMiddlewares.ResponseHandler,
			)

			s.Group("/", func(group *ghttp.RouterGroup) {
				// 首页HTML
				group.ALL("/", func(r *ghttp.Request) {
					r.Response.Write(g_index.IndexHTML)
				})
				group.ALL("/version", func(r *ghttp.Request) {
					r.Response.Write(binInfo.VersionString())
				})
				// 接口绑定
				r_hamster_router.BindController(group)
			})

			// 初始化
			if err := boot.Boot(); err != nil {
				glog.Fatal(ctx, "初始化任务失败: ", err)
			}

			s.Run()
			return nil
		},
	}
)
