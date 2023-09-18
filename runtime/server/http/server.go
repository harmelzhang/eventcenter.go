package http

import (
	"eventcenter-go/runtime/consts"
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/plugins/registry"
	"eventcenter-go/runtime/server"
	"eventcenter-go/runtime/server/http/handler"
	"eventcenter-go/runtime/server/http/router"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"log"
)

type HTTPServer struct {
	server *ghttp.Server
}

func New() server.CoreServer {
	return &HTTPServer{
		server: g.Server(),
	}
}

// Start 启动服务
func (s HTTPServer) Start() {
	registryPlugin := plugins.GetActivedPluginByType(plugins.TypeRegistry)
	if registryPlugin != nil {
		registryService := registryPlugin.(registry.Plugin).Service()
		err := registryService.Register(s.server.GetName(), s.server.GetListenedAddress(), consts.ProtocolHTTP)
		if err != nil {
			log.Panicf("register server err: %v", err)
		}
	}

	s.server.Use(handler.ErrorHandlerMiddleware)
	s.server.Use(ghttp.MiddlewareHandlerResponse)
	s.server.Group("/", func(group *ghttp.RouterGroup) {
		router.Router(group)
	})
	s.server.Run()
}

// Stop 停止服务
func (s HTTPServer) Stop() error {
	return s.server.Shutdown()
}
