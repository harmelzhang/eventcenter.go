package controller

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/plugins/connector"
	"eventcenter-go/runtime/plugins/storage"
)

var storagePlugin storage.StoragePlugin
var connectorPlugin connector.ConnectorPlugin

// RegisterStoragePlugin 注册存储插件
func RegisterStoragePlugin() {
	storagePlugin = plugins.GetActivedPluginByType(plugins.TypeStorage).(storage.StoragePlugin)
}

// RegisterConnectorPlugin 注册连接器插件
func RegisterConnectorPlugin() {
	connectorPlugin = plugins.GetActivedPluginByType(plugins.TypeConnector).(connector.ConnectorPlugin)
}
