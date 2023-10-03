// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package data_core

import (
	"context"

	"home-network-watcher/api/data_core/v1"
)

type IDataCoreV1 interface {
	GetXuiDataSSE(ctx context.Context, req *v1.GetXuiDataSSEReq) (res *v1.GetXuiDataSSERes, err error)
	GetNetworkDataSSE(ctx context.Context, req *v1.GetNetworkDataSSEReq) (res *v1.GetNetworkDataSSERes, err error)
	GetDockerMonitorSSE(ctx context.Context, req *v1.GetDockerMonitorSSEReq) (res *v1.GetDockerMonitorSSERes, err error)
	GetUptimeDataSSE(ctx context.Context, req *v1.GetUptimeDataSSEReq) (res *v1.GetUptimeDataSSERes, err error)
	GetHomeDataSSE(ctx context.Context, req *v1.GetHomeDataSSEReq) (res *v1.GetHomeDataSSERes, err error)
	GetTaskList(ctx context.Context, req *v1.GetTaskListReq) (res *v1.GetTaskListRes, err error)
	StopTask(ctx context.Context, req *v1.StopTaskReq) (res *v1.StopTaskRes, err error)
	RecoverTask(ctx context.Context, req *v1.RecoverTaskReq) (res *v1.RecoverTaskRes, err error)
}
