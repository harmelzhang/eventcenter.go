package admin

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/plugins/storage"
)

var storagePlugin storage.StoragePlugin

// RegisterStoragePlugin 注册存储插件
func RegisterStoragePlugin() {
	storagePlugin = plugins.GetActivedPluginByType(plugins.TypeStorage).(storage.StoragePlugin)
}
