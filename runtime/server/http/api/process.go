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
	TopicName string `p:"topic" v:"required#主题名称不能为空" dc:"主题名称"`
}

type TriggerRes struct {
	BaseRes
	EventId string `json:"event_id"`
}
