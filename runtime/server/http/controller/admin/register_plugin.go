package admin

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/plugin/storage"
)

var storagePlugin storage.StoragePlugin

// RegisterStoragePlugin 注册存储插件
func RegisterStoragePlugin() {
	storagePlugin = plugin.GetActivedPluginByType(plugin.TypeStorage).(storage.StoragePlugin)
}
