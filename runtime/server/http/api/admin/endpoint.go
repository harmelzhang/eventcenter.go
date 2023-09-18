package admin

import (
	"eventcenter-go/runtime/server/http/api"
	"github.com/gogf/gf/v2/frame/g"
)

// 查询终端

type QueryEndpointReq struct {
	g.Meta     `path:"/" method:"get" tags:"终端" summary:"查询终端"`
	TopicName  string `p:"topic" dc:"主题名称"`
	Type       string `p:"type" dc:"事件类型"`
	ServerName string `p:"server" dc:"服务名称"`
	Protocol   string `p:"protocol" dc:"协议名称"`
	Offset     int    `p:"offset" dc:"跳过多少条"`
	Limit      int    `p:"limit" dc:"取多少条"`
}

type QueryEndpointRes struct {
	api.PageRes
}

// 删除终端

type DeleteEndpointReq struct {
	g.Meta `path:"/" method:"delete" tags:"终端" summary:"删除终端"`
	Id     string `p:"id" v:"required#终端ID不能为空" dc:"ID"`
}

type DeleteEndpointRes struct {
	api.BaseRes
}

// 创建终端

type CreateEndpointReq struct {
	g.Meta     `path:"/" method:"post" tags:"终端" summary:"创建终端"`
	TopicName  string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
	Type       string `p:"type" v:"required#事件类型不能为空" dc:"事件类型"`
	ServerName string `p:"server" v:"required#服务名称不能为空" dc:"服务名称"`
	IsMicro    int    `p:"is_micro" v:"required#是否微服务不能为空" dc:"是否微服务"`
	Protocol   string `p:"protocol" v:"required#协议名称不能为空" dc:"协议名称"`
	Url        string `p:"url" v:"required#事件处理地址不能为空" dc:"事件处理地址"`
}

type CreateEndpointRes struct {
	api.BaseRes
}
