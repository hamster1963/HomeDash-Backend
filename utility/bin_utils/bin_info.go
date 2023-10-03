package binInfo

// 初始化为 unknown，如果编译时没有传入这些值，则为 unknown
var (
	GitTag         = "unknown"
	GitCommitLog   = "unknown"
	GitStatus      = "cleanly"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
)

// VersionString 返回版本信息
func VersionString() string {
	return "GitTag:" + GitTag + "\n" +
		"GitCommitLog:" + GitCommitLog + "\n" +
		"GitStatus:" + GitStatus + "\n" +
		"BuildTime:" + BuildTime + "\n" +
		"BuildGoVersion:" + BuildGoVersion + "\n"
}
