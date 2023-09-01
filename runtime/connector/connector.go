package connector

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// LifeCycle 生命周期方法
type LifeCycle interface {
	// Init 初始化
	Init() error

	// IsStarted 是否启动
	IsStarted() bool

	// IsStoped 是否停止
	IsStoped() bool

	// Start 启动服务
	Start() error

	// Stop 停止服务
	Stop() error
}

// Producer 生产者
type Producer interface {
	LifeCycle

	// Publish 发布事件
	Publish(ctx context.Context, event *cloudevents.Event) (err error)
}

// Consumer 消费者
type Consumer interface {
	LifeCycle

	// Subscribe 订阅
	Subscribe(topicName string) (err error)

	// Unsubscribe 取消订阅
	Unsubscribe(topicName string) (err error)

	// RegisterHandler 注册事件处理函数
	RegisterHandler(handler *EventHandler)
}

type EventHandler struct {
	Handler HandlerFunc
}

type HandlerFunc func(event *cloudevents.Event) (err error)
