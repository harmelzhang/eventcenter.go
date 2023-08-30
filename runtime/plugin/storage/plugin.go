package storage

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
)

type StoragePlugin interface {
	plugin.Plugin

	// TopicService 主题数据访问层
	TopicService() storage.TopicService
}
