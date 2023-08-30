package plugin

// Plugin 插件
type Plugin interface {
	// Type 插件类型
	Type() string

	// Init 初始化
	Init() error
}

// 注册的所有插件
// 类型: [插件类型]<[插件名称]插件>
var registedPlugins = make(map[string]map[string]Plugin)

// 激活的插件
// 类型: 插件名称
var activePlugins = make(map[string]string)

// Register 注册插件
func Register(name string, plugin Plugin) {
	plugins, isOK := registedPlugins[plugin.Type()]
	if !isOK {
		plugins = make(map[string]Plugin)
		registedPlugins[plugin.Type()] = plugins
	}
	plugins[name] = plugin
}

// ActivePlugin 激活插件
func ActivePlugin(typ, name string) {
	activePlugins[typ] = name
}

// GetPlugins 获取所有插件
func GetPlugins() map[string]map[string]Plugin {
	return registedPlugins
}

// Get 获取插件
func Get(typ string, name string) Plugin {
	return registedPlugins[typ][name]
}

// GetActivedPluginByType 根据类型获取激活的插件
func GetActivedPluginByType(typ string) Plugin {
	return Get(typ, activePlugins[typ])
}
