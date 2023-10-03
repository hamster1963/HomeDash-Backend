package g_functions

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

func ResErr(code int, errMsg ...interface{}) error {
	resMsg := ""
	errDetail := ""
	switch len(errMsg) {
	case 0:
		resMsg = "Unknown error reason"
	case 1:
		resMsg = gconv.String(errMsg[0])
	default:
		for _, v := range gconv.SliceStr(errMsg[:len(errMsg)-1]) {
			resMsg += v + ","
		}
		resMsg = resMsg[:len(resMsg)-1]
		errDetail = gconv.String(errMsg[len(errMsg)-1])
	}
	return gerror.NewCode(gcode.New(code, resMsg, nil), errDetail)
}

// SetDefaultHandler 替代默认的日志handler
func SetDefaultHandler() {
	glog.SetDefaultHandler(func(ctx context.Context, in *glog.HandlerInput) {
		m := map[string]interface{}{
			"stdout":            true,
			"writerColorEnable": false,
			"file":              "inside-{Y-m-d}.log",
			"path":              g.Config().MustGet(ctx, "server.logger.path", "log/").String(), // 此处必须重新设置，才可以实现db的log写入文件
		}
		_ = in.Logger.SetConfigWithMap(m)
		in.Next(ctx)
	})
}
