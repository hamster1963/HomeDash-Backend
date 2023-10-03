package boot

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"home-network-watcher/utility/ha_utils"
	"home-network-watcher/utility/network_utils"
	"home-network-watcher/utility/server_utils"
	"home-network-watcher/utility/uptime_utils"
)

// Boot
//
// @dc: 定时任务启动
// @return: error
func Boot() (err error) {
	_, err = gcron.AddOnce(context.TODO(), "@every 1s", func(ctx context.Context) {
		glog.Debug(context.Background(), "定时任务启动中...")
		if err := bootMethod(); err != nil {
			glog.Fatal(context.Background(), "定时任务启动失败: ", err)
		}
		glog.Debug(context.Background(), "定时任务启动成功")
	}, "开始启动定时任务")
	if err != nil {
		return err
	}

	_, err = gcron.AddOnce(context.TODO(), "@every 15s", func(ctx context.Context) {
		glog.Info(context.Background(), "定时任务测试中...")
		if err := bootCheck(); err != nil {
			glog.Fatal(context.Background(), "定时任务测试失败: ", err)
		} else {
			glog.Info(context.Background(), "定时任务测试成功")
		}
	}, "开始测试定时任务")
	if err != nil {
		return err
	}

	return nil
}

// bootCheck
//
// @Description: 定时任务测试
// @return error
func bootCheck() (err error) {
	return nil
}

// bootMethod
// @Description: 定时任务启动
// @return error
func bootMethod() (err error) {
	var ctx = gctx.New()

	glog.Notice(ctx, "开始获取科学上网网速")
	_, err = gcron.AddSingleton(ctx, "@every 1s", func(ctx context.Context) {
		err = network_utils.ProxyNetwork.GetProxyNetwork()
		if err != nil {
			glog.Warning(ctx, "获取代理速度"+err.Error())
		}
	}, "获取代理速度")
	if err != nil {
		return err
	}

	glog.Notice(ctx, "开始获取家庭路由器网速")
	_, err = gcron.AddSingleton(ctx, "@every 1s", func(ctx context.Context) {
		err = network_utils.NetworkUtils.GetHomeNetwork()
		if err != nil {
			glog.Warning(ctx, "获取家庭路由器速度"+err.Error())
		}
	}, "获取家庭路由器速度")
	if err != nil {
		return err
	}

	glog.Notice(ctx, "开始获取当前代理节点信息")
	_, err = gcron.AddSingleton(ctx, "@every 5s", func(ctx context.Context) {
		err = network_utils.NodeUtils.GetNodeInfo()
		if err != nil {
			glog.Warning(ctx, "获取当前代理节点信息"+err.Error())
		}
	}, "获取当前代理节点信息")
	if err != nil {
		return err
	}

	// 进行第一次机场信息缓存
	err = network_utils.ProxyProvider.GetSubscribeInfo()
	if err != nil {
		glog.Warning(ctx, "获取机场订阅信息失败"+err.Error())
	}

	glog.Notice(ctx, "开始获取机场订阅信息")
	_, err = gcron.AddSingleton(ctx, "@every 10m", func(ctx context.Context) {
		err = network_utils.ProxyProvider.GetSubscribeInfo()
		if err != nil {
			glog.Warning(ctx, "获取机场订阅信息失败"+err.Error())
		}
	}, "开始获取机场订阅信息")
	if err != nil {
		return err
	}

	// 获取服务器信息
	glog.Notice(ctx, "开始获取服务器信息")
	_, err = gcron.AddSingleton(ctx, "@every 1s", func(ctx context.Context) {
		err = server_utils.Nezha.GetAllServerInfo(ctx)
		if err != nil {
			glog.Warning(ctx, "获取服务器信息失败"+err.Error())
		}
	}, "开始获取服务器信息")
	if err != nil {
		return err
	}

	// 获取 xui 信息
	glog.Notice(ctx, "开始获取 xui 信息")
	_, err = gcron.AddSingleton(ctx, "@every 5s", func(ctx context.Context) {
		err = network_utils.ProxyNetwork.GetXuiUserList()
		if err != nil {
			glog.Warning(ctx, "获取 xui 信息失败"+err.Error())
		}
	}, "开始获取 xui 信息")
	if err != nil {
		return err
	}

	// 获取服务概览信息
	glog.Notice(ctx, "开始获取服务概览信息")
	_, err = gcron.AddSingleton(ctx, "@every 10s", func(ctx context.Context) {
		err = server_utils.ServerCron.CronGetDockerAndMonitor(ctx)
		if err != nil {
			glog.Warning(ctx, "获取服务概览信息失败"+err.Error())
		}
	}, "开始获取服务概览信息")
	if err != nil {
		return err
	}

	// 获取 UptimeKuma 信息
	glog.Notice(ctx, "开始获取 UptimeKuma 信息")
	_, err = gcron.AddSingleton(ctx, "@every 10s", func(ctx context.Context) {
		err = uptime_utils.Kuma.GetUptimeData(ctx)
		if err != nil {
			glog.Warning(ctx, "获取 UptimeKuma 信息失败"+err.Error())
		}
	}, "开始获取 UptimeKuma 信息")
	if err != nil {
		return err
	}

	// 获取智能家居信息
	glog.Notice(ctx, "开始获取智能家居信息")
	_, err = gcron.AddSingleton(ctx, "@every 1s", func(ctx context.Context) {
		err = ha_utils.HaUtils.GetEntitiesInfo(ctx)
		if err != nil {
			glog.Warning(ctx, "获取智能家居信息失败"+err.Error())
		}
	}, "开始获取智能家居信息")
	if err != nil {
		return err
	}

	return nil
}
