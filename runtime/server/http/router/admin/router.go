package admin

import (
	"eventcenter-go/runtime/server/http/controller/admin"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Router(group *ghttp.RouterGroup) {
	group.Group("/topic", func(group *ghttp.RouterGroup) {
		group.Bind(admin.TopicController)
	})
	group.Group("/event", func(group *ghttp.RouterGroup) {
		group.Bind(admin.EventController)
	})
	group.Group("/endpoint", func(group *ghttp.RouterGroup) {
		group.Bind(admin.EndpointController)
	})
}
