package connector

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugin"
)

type ConnectorPlugin interface {
	plugin.Plugin

	// Producer 获取生产者
	Producer() (connector.Producer, error)

	// Consumer 获取消费者
	Consumer() (connector.Consumer, error)
}
