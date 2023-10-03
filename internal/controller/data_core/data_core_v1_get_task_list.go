package data_core

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"

	"home-network-watcher/api/data_core/v1"
)

// GetTaskList 获取任务列表
func (c *ControllerV1) GetTaskList(_ context.Context, _ *v1.GetTaskListReq) (res *v1.GetTaskListRes, err error) {
	var taskList []map[string]interface{}
	for _, entry := range gcron.Entries() {
		taskList = append(taskList, map[string]interface{}{
			"name":   entry.Name,
			"status": entry.Status(),
			"time":   entry.Time,
		})
	}
	res = &v1.GetTaskListRes{
		TaskList: taskList,
	}
	return res, nil
}
