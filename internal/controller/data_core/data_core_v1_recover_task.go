package data_core

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"home-network-watcher/internal/global/g_functions"

	"home-network-watcher/api/data_core/v1"
)

// RecoverTask 恢复任务
func (c *ControllerV1) RecoverTask(_ context.Context, req *v1.RecoverTaskReq) (res *v1.RecoverTaskRes, err error) {
	var taskMap = make(map[string]int)
	for _, entry := range gcron.Entries() {
		taskMap[entry.Name] = entry.Status()
	}
	if _, ok := taskMap[req.Name]; !ok {
		return nil, g_functions.ResErr(400, "任务不存在")
	}
	if taskMap[req.Name] == 1 || taskMap[req.Name] == 0 {
		return nil, g_functions.ResErr(400, "任务已启动")
	}
	gcron.Start(req.Name)
	return
}
