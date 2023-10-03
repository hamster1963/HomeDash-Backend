package data_core

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"home-network-watcher/internal/global/g_functions"

	"home-network-watcher/api/data_core/v1"
)

// StopTask 停止任务
func (c *ControllerV1) StopTask(_ context.Context, req *v1.StopTaskReq) (res *v1.StopTaskRes, err error) {
	var taskMap = make(map[string]string)
	for _, entry := range gcron.Entries() {
		taskMap[entry.Name] = entry.Name
	}
	if _, ok := taskMap[req.Name]; !ok {
		return nil, g_functions.ResErr(400, "任务不存在")
	}
	gcron.Stop(req.Name)
	return
}
