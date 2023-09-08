package runtime

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/consts"
	"eventcenter-go/runtime/plugins"
	connectorPlugin "eventcenter-go/runtime/plugins/connector"
	"eventcenter-go/runtime/server"
	"eventcenter-go/runtime/server/grpc"
	"eventcenter-go/runtime/server/http"
	"eventcenter-go/runtime/server/http/controller"
	"eventcenter-go/runtime/server/http/controller/admin"
	"eventcenter-go/runtime/server/tcp"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"log"
)

import (
	// 加载数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	// 加载存储插件
	_ "eventcenter-go/runtime/plugins/storage/database"
	_ "eventcenter-go/runtime/plugins/storage/mongodb"
	_ "eventcenter-go/runtime/plugins/storage/redis"
	_ "eventcenter-go/runtime/plugins/storage/standalone"
	// 加载连接器插件
	_ "eventcenter-go/runtime/plugins/connector/redis"
	_ "eventcenter-go/runtime/plugins/connector/standalone"
)

// Start 启动所有服务
func Start() error {
	var initSuccessed bool
	var servers []server.CoreServer

	defer func() {
		// 初始化失败
		if !initSuccessed && len(servers) > 0 {
			for _, srv := range servers {
				_ = srv.Stop()
			}
		}
	}()

	httpServer := http.New()
	servers = append(servers, httpServer)

	config := g.Cfg().MustData(gctx.New())

	// TCP Server
	if _, isOK := config[consts.ConfigTcp]; isOK {
		tcpServer := tcp.New()
		servers = append(servers, tcpServer)
	}

	// gRPC Server
	if _, isOK := config[consts.ConfigGrpc]; isOK {
		grpcServer := grpc.New()
		servers = append(servers, grpcServer)
	}

	// 这里可以换成 fix 包，通过依赖注入框架实现生命周期管理
	wait := make(chan bool)
	for _, srv := range servers {
		go srv.Start()
	}
	<-wait

	return nil
}

// LoadPlugins 加载插件
func LoadPlugins() error {
	config := g.Cfg().MustData(gctx.New())
	if _, isOK := config[consts.ConfigPlugins]; isOK {
		ctx := gctx.New()
		cfg, err := g.Cfg().Get(ctx, "plugins")
		if err != nil {
			return err
		}

		// 循环加载插件
		for key, value := range cfg.MapStrVar() {
			err = loadPlugins(key, value.MapStrVar())
			if err != nil {
				return err
			}
		}

		// 注册插件
		registerPlugins()
	}
	return nil
}

// 加载插件
func loadPlugins(pluginType string, config map[string]*gvar.Var) error {
	// 激活插件
	activePluginName := getActivePluginName(pluginType, config)
	plugins.ActivePlugin(pluginType, activePluginName)

	// 初始化插件
	err := initPlugin(pluginType, activePluginName, config)
	if err != nil {
		return err
	}

	if pluginType == plugins.TypeConnector {
		plugin := plugins.GetActivedPluginByType(plugins.TypeConnector).(connectorPlugin.Plugin)
		consumer, err := plugin.Consumer()
		if err != nil {
			return err
		}
		consumer.RegisterHandler(connector.NewEventHandler())
	}

	return nil
}

// 注册插件
func registerPlugins() {
	// 存储
	controller.RegisterStoragePlugin()
	admin.RegisterStoragePlugin()
	// 连接器
	controller.RegisterConnectorPlugin()
	admin.RegisterConnectorPlugin()
}

// 获取激活插件名
func getActivePluginName(pluginType string, config map[string]*gvar.Var) string {
	active, isOK := config[plugins.NameActive]

	activePluginName := plugins.NameStorageStandalone
	if pluginType == plugins.TypeConnector {
		activePluginName = plugins.NameConnectorStandalone
	}

	if isOK {
		name := active.String()
		if _, has := config[name]; !has {
			if name != activePluginName {
				log.Printf("[%s] not found [%s] , use default config [%s]", pluginType, name, activePluginName)
			}
		} else {
			activePluginName = name
		}
	}
	return activePluginName
}

// 初始化插件
func initPlugin(pluginType, activePluginName string, config map[string]*gvar.Var) error {
	// 初始化插件
	p := plugins.Get(pluginType, activePluginName)
	cfg := make(map[string]*gvar.Var)
	if activePluginName != plugins.NameStorageStandalone {
		cfg = config[activePluginName].MapStrVar()
	}

	err := p.Init(cfg)
	if err != nil {
		return err
	}

	return nil
}
