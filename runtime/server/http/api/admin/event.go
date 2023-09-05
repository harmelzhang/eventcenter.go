package admin

import (
	"eventcenter-go/runtime/server/http/api"
	"github.com/gogf/gf/v2/frame/g"
)

// 查询事件

type QueryEventReq struct {
	g.Meta    `path:"/" method:"get" tags:"事件" summary:"查询事件"`
	Source    string `p:"source" dc:"事件源"`
	TopicName string `p:"topic" dc:"主题名称"`
	Type      string `p:"type" dc:"事件类型"`
	Offset    int    `p:"offset" dc:"跳过多少条"`
	Limit     int    `p:"limit" dc:"取多少条"`
}

type QueryEventRes struct {
	api.PageRes
}

// 删除事件

type DeleteEventReq struct {
	g.Meta `path:"/" method:"delete" tags:"事件" summary:"删除事件"`
	Id     string `p:"id" v:"required#事件ID不能为空" dc:"ID"`
}

type DeleteEventRes struct {
	api.BaseRes
}

// 创建事件

type CreateEventReq struct {
	g.Meta    `path:"/" method:"post" tags:"事件" summary:"创建事件"`
	Source    string `p:"source" v:"required#事件源不能为空" dc:"事件源名称"`
	TopicName string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
	Type      string `p:"type" v:"required#事件类型不能为空" dc:"事件类型"`
	Data      string `p:"data" v:"required#事件上下文信息不能为空" dc:"事件数据"`
}

type CreateEventRes struct {
	api.BaseRes
}
