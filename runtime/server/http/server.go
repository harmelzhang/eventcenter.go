package http

import (
	"eventcenter-go/runtime/server"
	"eventcenter-go/runtime/server/http/handler"
	"eventcenter-go/runtime/server/http/router"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HTTPServer struct {
	server *ghttp.Server
}

func New() server.CoreServer {
	return &HTTPServer{
		server: g.Server(),
	}
}

func (s HTTPServer) Start() {
	s.server.Use(handler.ErrorHandlerMiddleware)
	s.server.Use(ghttp.MiddlewareHandlerResponse)
	s.server.Group("/", func(group *ghttp.RouterGroup) {
		router.Router(group)
	})
	s.server.Run()
}

func (s HTTPServer) Stop() error {
	return s.server.Shutdown()
}
