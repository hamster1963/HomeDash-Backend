package docker_utils

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"home-network-watcher/manifest"
)

type uDockerUtils struct{}

var DockerUtils = &uDockerUtils{}

type dockerStatus struct {
	ServerCount int
	ErrorServer int
	DockerCount int
	ErrorDocker int
}

// GetDockerStatus
//
//	@dc: 获取docker状态
//	@author: hamster   @date:2023/9/20 15:33:19
func (u *uDockerUtils) GetDockerStatus(ctx context.Context) (status *dockerStatus, err error) {
	status = &dockerStatus{}
	response, err := g.Client().SetHeaderMap(manifest.DockerAuthMap).Get(context.Background(), manifest.DockerApiUrl)
	if err != nil {
		return
	}
	defer func(response *gclient.Response) {
		err = response.Close()
		if err != nil {
			glog.Warning(ctx, "关闭DockerStatus", err)
		}
	}(response)
	endpointList := gconv.SliceAny(response.ReadAllString())
	status.ServerCount = len(endpointList)
	for _, endpoint := range endpointList {
		endpointJson := gjson.New(endpoint)
		// 获取全部容器数量
		serverRunningDockerCount := endpointJson.Get("Snapshots.0.RunningContainerCount").Int()
		serverStoppedDockerCount := endpointJson.Get("Snapshots.0.StoppedContainerCount").Int()
		status.DockerCount += serverRunningDockerCount + serverStoppedDockerCount
		// 判断服务器状态
		if endpointJson.Get("Status").Int() != 1 {
			status.ErrorServer++
			status.ErrorDocker += serverRunningDockerCount + serverStoppedDockerCount
		} else {
			status.ErrorDocker += serverStoppedDockerCount
		}
	}
	return status, nil
}
