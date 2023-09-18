package registry

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/registry"
)

type Plugin interface {
	plugins.Plugin

	// Service 注册服务
	Service() registry.Service
}
