package handler

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// ErrorHandlerMiddleware 异常处理中间件
func ErrorHandlerMiddleware(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		r.Response.ClearBuffer()
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    -1,
			Message: err.Error(),
			Data:    nil,
		})
	}
}
