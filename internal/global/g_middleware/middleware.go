package g_middleware

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/golang-jwt/jwt/v5"
	"home-network-watcher/manifest"
	"sync"
	"time"
)

type sMiddleware struct{} // 创建结构体
var (
	once         = &sync.Once{}    // 创建互锁
	s            *sMiddleware      // 创建指针
	SMiddlewares = newMiddleware() // 对外暴露
)

// defaultHandlerResponse 返回结构体
type defaultHandlerResponse struct {
	Status  int         `json:"status"  dc:"Error code"` // 可业务需要可以更改json字段
	Message string      `json:"msg"     dc:"Error message"`
	Data    interface{} `json:"data"    dc:"Result data for certain request according API definition"`
}

// newMiddleware 单例中间件
func newMiddleware() *sMiddleware {
	once.Do(func() {
		s = &sMiddleware{}
	})
	return s
}

func (s *sMiddleware) MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// resWriteJson 返回json输出
func (s *sMiddleware) resWriteJson(r *ghttp.Request, in defaultHandlerResponse) {
	r.Response.ClearBuffer()
	r.Response.WriteJson(defaultHandlerResponse{
		Status:  in.Status,
		Message: in.Message,
		Data:    in.Data,
	})
}

// ErrorsStatus 服务器错误码处理
func (s *sMiddleware) ErrorsStatus(server *ghttp.Server) {
	server.BindStatusHandlerByMap(map[int]ghttp.HandlerFunc{
		500: func(r *ghttp.Request) {
			s.resWriteJson(r, defaultHandlerResponse{
				Status:  500,
				Message: "Error 500,Internal Server Error",
				Data:    "",
			})
		},
		404: func(r *ghttp.Request) {
			s.resWriteJson(r, defaultHandlerResponse{
				Status:  404,
				Message: "Error 404,Not Found",
				Data:    "",
			})
		},
	})
}

// JWTAuth 鉴权中间件
func (s *sMiddleware) JWTAuth(r *ghttp.Request) {
	JWTString := r.GetHeader("Authorization")
	if JWTString == "" {
		s.resWriteJson(r, defaultHandlerResponse{
			Status:  401,
			Message: "JWT为空",
			Data:    nil,
		})
		return
	}

	// 验证JWT
	token, err := jwt.Parse(JWTString, func(token *jwt.Token) (interface{}, error) {
		return manifest.JWTKey, nil
	})
	if err != nil {
		s.resWriteJson(r, defaultHandlerResponse{
			Status:  401,
			Message: "JWT验证失败",
			Data:    nil,
		})
		return
	}

	if !token.Valid {
		s.resWriteJson(r, defaultHandlerResponse{
			Status:  401,
			Message: "JWT验证失败",
			Data:    nil,
		})
		return
	}
	// 验证是否过期
	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		s.resWriteJson(r, defaultHandlerResponse{
			Status:  401,
			Message: "JWT验证失败",
			Data:    nil,
		})
		return
	}
	if expirationTime.Before(time.Now()) {
		s.resWriteJson(r, defaultHandlerResponse{
			Status:  401,
			Message: "JWT已过期",
			Data:    nil,
		})
		return
	}

	glog.Info(context.TODO(), "JWT验证通过")
	// 验证通过，设置用户信息
	audience, err := token.Claims.GetAudience()
	if err != nil {
		s.resWriteJson(r, defaultHandlerResponse{
			Status:  401,
			Message: "JWT验证失败",
			Data:    nil,
		})
		return
	}

	r.SetCtxVar("user_id", audience[0])
	r.Middleware.Next()
}

// ResponseHandler is the middleware handling handler response object and its error. 中间件处理处理程序响应对象及其错误
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	//defer action_log.NewActionLog().LogAdd(r) // 日志钩子
	r.Middleware.Next()

	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}

	// 定义参数和类型.方便调用
	var (
		code    int
		msg     string
		resData = r.GetHandlerResponse() // 检索并返回处理程序响应对象及其错误
		err     = r.GetError()           // 返回请求过程中发生的错误。如果没有错误，则返回 nil
		errCode = gerror.Code(err)       // 将错误信息通过接口解析处理
	)
	code = 200      // 默认返回状态码信息码
	msg = "success" // 默认成功返回响应体文本信息
	if err != nil {
		code = 400 // 默认返回状态码信息码
		switch errCode.Code() {
		case -1:
			msg = "Nonstandard error return"
			resData = "I won't expose it to you. Go back and change it!"
		case 500:
			r.Response.Writer.Status = 500
			return
		default:
			msg = errCode.Message()
			resData = err.Error()
			r.Response.Writer.Status = 400
			if msg == "" || msg == "Unknown error reason" {
				msg = "Nonstandard error return"
				resData = "I won't expose it to you. Go back and change it!"
			}
			if resData == msg {
				resData = ""
			}
		}
	}
	// 如果是空res返回，则resData返回空字符串
	if g.IsNil(resData) {
		resData = ""
	}
	s.resWriteJson(r, defaultHandlerResponse{
		Status:  code,
		Message: msg,
		Data:    resData,
	})
}
