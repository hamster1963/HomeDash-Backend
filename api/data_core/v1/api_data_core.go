package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetXuiDataSSEReq 获取 xui 网络数据 Req请求
type GetXuiDataSSEReq struct {
	g.Meta `method:"get" tags:"xui" summary:"获取 xui 网络数据" dc:"获取 xui 网络数据"`
}

// GetXuiDataSSERes 获取 xui 网络数据 Res返回
type GetXuiDataSSERes struct {
	XuiData interface{} `json:"xuiData" dc:"xui数据"`
}

// GetNetworkDataSSEReq 获取网络信息-SSE Req请求
type GetNetworkDataSSEReq struct {
	g.Meta `method:"get" tags:"家庭网络" summary:"获取网络信息-SSE" dc:"获取网络信息-SSE"`
}

// GetNetworkDataSSERes 获取网络信息-SSE Res返回
type GetNetworkDataSSERes struct {
	NodeInfo     interface{} `json:"nodeInfo" dc:"节点信息"`
	HomeNetwork  interface{} `json:"homeNetwork" dc:"家庭网络"`
	ProxyNetwork interface{} `json:"proxyNetwork" dc:"科学上网"`
	CoffeeInfo   interface{} `json:"coffeeInfo" dc:"coffee代理信息"`
	ServerInfo   interface{} `json:"serverInfo" dc:"服务器信息"`
}

// GetDockerMonitorSSEReq 获取 Docker 监控数据 Req请求
type GetDockerMonitorSSEReq struct {
	g.Meta `method:"get" tags:"docker" summary:"获取 Docker 监控数据" dc:"获取 Docker 监控数据"`
}

// GetDockerMonitorSSERes 获取 Docker 监控数据 Res返回
type GetDockerMonitorSSERes struct {
	DockerData interface{} `json:"dockerData" dc:"Docker数据"`
}

// GetUptimeDataSSEReq 获取 uptime 数据 Req请求
type GetUptimeDataSSEReq struct {
	g.Meta `method:"get" tags:"uptime" summary:"获取 uptime 数据" dc:"获取 uptime 数据"`
}

// GetUptimeDataSSERes 获取 uptime 数据 Res返回
type GetUptimeDataSSERes struct {
	UptimeData interface{} `json:"uptimeData" dc:"uptime数据"`
}

// GetHomeDataSSEReq 获取智能家居数据 Req请求
type GetHomeDataSSEReq struct {
	g.Meta `method:"get" tags:"home" summary:"获取智能家居数据" dc:"获取智能家居数据"`
}

// GetHomeDataSSERes 获取智能家居数据 Res返回
type GetHomeDataSSERes struct {
	HomeData interface{} `json:"homeData" dc:"智能家居数据"`
}

// GetTaskListReq 获取定时任务列表 Req请求
type GetTaskListReq struct {
	g.Meta `method:"get" tags:"任务" summary:"获取定时任务列表" dc:"获取定时任务列表"`
}

// GetTaskListRes 获取定时任务列表 Res返回
type GetTaskListRes struct {
	TaskList interface{} `json:"taskList" dc:"定时任务列表"`
}

// StopTaskReq 停止定时任务 Req请求
type StopTaskReq struct {
	g.Meta `method:"post" tags:"任务" summary:"停止定时任务" dc:"停止定时任务"`
	Name   string `json:"name" dc:"任务名称" v:"required #请输入 任务名称"`
}

// StopTaskRes 停止定时任务 Res返回
type StopTaskRes struct {
}

// RecoverTaskReq 恢复定时任务 Req请求
type RecoverTaskReq struct {
	g.Meta `method:"post" tags:"任务" summary:"恢复定时任务" dc:"恢复定时任务"`
	Name   string `json:"name" dc:"任务名称" v:"required #请输入 任务名称"`
}

// RecoverTaskRes 恢复定时任务 Res返回
type RecoverTaskRes struct {
}
