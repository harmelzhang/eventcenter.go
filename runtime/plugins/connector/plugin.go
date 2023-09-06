package connector

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugins"
)

type Plugin interface {
	plugins.Plugin

	// Producer 获取生产者
	Producer() (connector.Producer, error)

	// Consumer 获取消费者
	Consumer() (connector.Consumer, error)
}
