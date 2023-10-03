package server_utils

import (
	"context"
	"github.com/gogf/gf/v2/os/gcache"
	"home-network-watcher/manifest"
	"home-network-watcher/utility/docker_utils"
	"home-network-watcher/utility/uptime_utils"
	"time"
)

type uServerCron struct{}

var ServerCron = &uServerCron{}

// CronGetDockerAndMonitor
//
//	@dc: 获取 Docker 状态与监控信息
//	@author: hamster   @date:2023/9/20 16:20:11
func (u *uServerCron) CronGetDockerAndMonitor(ctx context.Context) (err error) {
	dockerStatus, _ := docker_utils.DockerUtils.GetDockerStatus(ctx)
	serverCount, errorCount, _ := uptime_utils.Kuma.GetMonitorStatus(ctx)
	// 缓存数据
	dockerMonitor := map[string]interface{}{
		"dockerStatus": dockerStatus,
		"serverCount":  serverCount,
		"errorCount":   errorCount,
	}
	err = gcache.Set(ctx, manifest.DockerMonitorCacheKey, dockerMonitor, 1*time.Hour)
	if err != nil {
		return err
	}
	return
}
