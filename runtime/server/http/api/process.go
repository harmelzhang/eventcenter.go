package api

import "github.com/gogf/gf/v2/frame/g"

// 订阅事件

type SubscribeReq struct {
	g.Meta     `path:"/subscribe" method:"post" tags:"事件处理" summary:"订阅事件"`
	TopicName  string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
	Type       string `p:"type" v:"required#事件类型不能为空" dc:"事件类型"`
	ServerName string `p:"server" v:"required#服务名称不能为空" dc:"服务名称"`
	IsMicro    int    `p:"is_micro" v:"required#是否微服务不能为空" dc:"是否微服务"`
	Protocol   string `p:"protocol" v:"required#协议名称不能为空" dc:"协议名称"`
	Url        string `p:"url" v:"required#事件处理地址不能为空" dc:"事件处理地址"`
}

type SubscribeRes struct {
	BaseRes
}

// 取消订阅

type UnsubscribeReq struct {
	g.Meta     `path:"/unsubscribe" method:"post" tags:"事件处理" summary:"取消订阅"`
	TopicName  string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
	Type       string `p:"type" v:"required#事件类型不能为空" dc:"事件类型"`
	ServerName string `p:"server" v:"required#服务名称不能为空" dc:"服务名称"`
	Protocol   string `p:"protocol" v:"required#协议名称不能为空" dc:"协议名称"`
}

type UnsubscribeRes struct {
	BaseRes
}

// 触发事件

type TriggerReq struct {
	g.Meta    `path:"/trigger" method:"post" tags:"事件处理" summary:"触发事件"`
	Source    string `p:"source" v:"required#事件源不能为空" dc:"事件源名称"`
	TopicName string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
	Type      string `p:"type" v:"required#事件类型不能为空" dc:"事件类型"`
	Data      string `p:"data" v:"required#事件上下文信息不能为空" dc:"事件数据"`
}

type TriggerRes struct {
	BaseRes
	EventId string `json:"event_id"`
}
