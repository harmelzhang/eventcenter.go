package storage

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
)

type StoragePlugin interface {
	plugins.Plugin

	// TopicService 主题数据访问层
	TopicService() storage.TopicService

	// EndpointService 终端数据访问层
	EndpointService() storage.EndpointService

	// EventService 事件数据访问层
	EventService() storage.EventService
}
