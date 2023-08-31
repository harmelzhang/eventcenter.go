package api

import "github.com/gogf/gf/v2/frame/g"

// 订阅事件

type SubscribeReq struct {
	g.Meta    `path:"/subscribe" method:"post" tags:"事件处理" summary:"订阅事件"`
	Url       string `p:"url" v:"required|url#事件处理地址不能为空|事件处理地址格式错误" dc:"事件处理地址"`
	TopicName string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
	//Type      string `p:"type" v:"in:SYNC,ASYNC#事件执行类型" dc:"事件执行类型（同步、异步）"`
}

type SubscribeRes struct {
	BaseRes
}

// 触发事件

type TriggerReq struct {
	g.Meta    `path:"/trigger" method:"post" tags:"事件处理" summary:"触发事件"`
	Source    string `p:"source" v:"required#事件源不能为空" dc:"事件源名称"`
	TopicName string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
	Data      string `p:"data" v:"required#事件上下文信息不能为空" dc:"事件数据"`
}

type TriggerRes struct {
	BaseRes
	EventId string `json:"event_id"`
}
