package router

import (
	"eventcenter-go/runtime/server/http/controller"
	"eventcenter-go/runtime/server/http/router/admin"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Router(group *ghttp.RouterGroup) {
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(controller.ProcessController)
	})
	group.Group("/admin", func(group *ghttp.RouterGroup) {
		admin.Router(group)
	})
}
