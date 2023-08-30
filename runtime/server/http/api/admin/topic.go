package admin

import (
	"eventcenter-go/runtime/server/http/api"
	"github.com/gogf/gf/v2/frame/g"
)

// 创建主题

type CreateTopicReq struct {
	g.Meta `path:"/" method:"post" tags:"主题" summary:"创建主题"`
	Name   string `p:"name" v:"required#主题名不能为空" dc:"名称"`
}

type CreateTopicRes struct {
	api.BaseRes
}

// 查询主题

type QueryTopicReq struct {
	g.Meta `path:"/" method:"get" tags:"主题" summary:"主题列表"`
}

type QueryTopicRes struct {
	api.PageRes
}

// 删除主题

type DeleteTopicReq struct {
	g.Meta `path:"/" method:"delete" tags:"主题" summary:"删除主题"`
	Id     string `p:"id" v:"required#ID不能为空" dc:"ID"`
}

type DeleteTopicRes struct {
	api.BaseRes
}
