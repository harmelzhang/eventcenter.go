package admin

import (
	"eventcenter-go/runtime/server/http/controller/admin"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Router(group *ghttp.RouterGroup) {
	group.Group("/topic", func(group *ghttp.RouterGroup) {
		group.Bind(admin.TopicController)
	})
}
