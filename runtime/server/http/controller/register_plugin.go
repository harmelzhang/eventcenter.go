package controller

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/plugin/connector"
	"eventcenter-go/runtime/plugin/storage"
)

var storagePlugin storage.StoragePlugin
var connectorPlugin connector.ConnectorPlugin

// RegisterStoragePlugin 注册存储插件
func RegisterStoragePlugin() {
	storagePlugin = plugin.GetActivedPluginByType(plugin.TypeStorage).(storage.StoragePlugin)
}

// RegisterConnectorPlugin 注册连接器插件
func RegisterConnectorPlugin() {
	connectorPlugin = plugin.GetActivedPluginByType(plugin.TypeConnector).(connector.ConnectorPlugin)
}
