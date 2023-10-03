package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"home-network-watcher/internal/cmd"
	_ "home-network-watcher/internal/logic"
	_ "home-network-watcher/internal/packed"
	binInfo "home-network-watcher/utility/bin_utils"
)

var (
	GitTag         = "unknown"
	GitCommitLog   = "unknown"
	GitStatus      = "cleanly"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
)

func main() {
	// 注入编译时的信息
	binInfo.GitTag = GitTag
	binInfo.GitCommitLog = GitCommitLog
	binInfo.GitStatus = GitStatus
	binInfo.BuildTime = BuildTime
	binInfo.BuildGoVersion = BuildGoVersion
	cmd.Main.Run(gctx.New())
}
