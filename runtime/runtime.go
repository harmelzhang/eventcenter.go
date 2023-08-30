package runtime

import (
	"eventcenter-go/runtime/consts"
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/plugin/storage/mongodb"
	"eventcenter-go/runtime/server"
	"eventcenter-go/runtime/server/grpc"
	"eventcenter-go/runtime/server/http"
	"eventcenter-go/runtime/server/http/controller"
	"eventcenter-go/runtime/server/http/controller/admin"
	"eventcenter-go/runtime/server/tcp"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

import (
	// 加载数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	// 加载存储插件
	_ "eventcenter-go/runtime/plugin/storage/database"
	_ "eventcenter-go/runtime/plugin/storage/mongodb"
	_ "eventcenter-go/runtime/plugin/storage/redis"
	_ "eventcenter-go/runtime/plugin/storage/standalone"
	// 加载连接器插件
	_ "eventcenter-go/runtime/plugin/connector/standalone"
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
			if key == plugin.TypeStorage {
				err = loadStoragePlugins(value)
				if err != nil {
					return err
				}
			} else if key == plugin.TypeConnector {
				err = loadConnectorPlugins(value)
				if err != nil {
					return err
				}
			}
		}

		// 注册插件
		registerPlugins()
	}
	return nil
}

// 加载存储插件
func loadStoragePlugins(cfgVar *gvar.Var) error {
	// 加载配置
	config := cfgVar.MapStrVar()

	// 循环注册插件
	for key, value := range config {
		if key == plugin.NameActive {
			continue
		}

		configInfo := value.MapStrVar()

		if key == plugin.NameStorageRedis {
			gredis.SetConfig(&gredis.Config{
				Address:         configInfo["address"].String(),
				Pass:            configInfo["password"].String(),
				Db:              configInfo["db"].Int(),
				User:            configInfo["user"].String(),
				MinIdle:         configInfo["minIdle"].Int(),
				MaxIdle:         configInfo["maxIdle"].Int(),
				MaxActive:       configInfo["maxActive"].Int(),
				MaxConnLifetime: configInfo["maxConnLifetime"].Duration() * time.Second,
				IdleTimeout:     configInfo["idleTimeout"].Duration() * time.Second,
				WaitTimeout:     configInfo["waitTimeout"].Duration() * time.Second,
				DialTimeout:     configInfo["dialTimeout"].Duration() * time.Second,
				ReadTimeout:     configInfo["readTimeout"].Duration() * time.Second,
				WriteTimeout:    configInfo["writeTimeout"].Duration() * time.Second,
				MasterName:      configInfo["masterName"].String(),
				TLS:             configInfo["tls"].Bool(),
				TLSSkipVerify:   configInfo["tlsSkipVerify"].Bool(),
				SlaveOnly:       configInfo["slaveOnly"].Bool(),
			}, plugin.TypeStorage)
		} else if key == plugin.NameStorageDB {
			gdb.SetConfig(gdb.Config{
				plugin.TypeStorage: gdb.ConfigGroup{
					gdb.ConfigNode{
						Host:             configInfo["host"].String(),
						Port:             configInfo["port"].String(),
						User:             configInfo["user"].String(),
						Pass:             configInfo["password"].String(),
						Name:             configInfo["name"].String(),
						Type:             configInfo["type"].String(),
						Link:             configInfo["link"].String(),
						Extra:            configInfo["extra"].String(),
						Role:             configInfo["role"].String(),
						Debug:            configInfo["debug"].Bool(),
						Charset:          configInfo["charset"].String(),
						Prefix:           configInfo["prefix"].String(),
						Weight:           configInfo["weight"].Int(),
						MaxIdleConnCount: configInfo["maxIdle"].Int(),
						MaxOpenConnCount: configInfo["maxOpen"].Int(),
						MaxConnLifeTime:  configInfo["maxLifetime"].Duration() * time.Second,
					},
				},
			})
		} else if key == plugin.NameStorageMongodb {
			clientOptions := options.Client().ApplyURI(configInfo["uri"].String())
			ctx := gctx.New()
			conn, err := mongo.Connect(ctx, clientOptions)
			poolSize := configInfo["poolSize"].Uint64()
			if poolSize > 0 {
				clientOptions.SetMaxPoolSize(poolSize)
			}
			if err != nil {
				return err
			}
			err = conn.Ping(ctx, nil)
			if err != nil {
				return err
			}
			mongodb.InitDB(conn.Database(configInfo["database"].String()))
		}
	}

	active, isOK := config[plugin.NameActive]
	activePluginName := plugin.NameStorageStandalone
	if isOK {
		name := active.String()
		if _, has := config[name]; !has {
			if name != plugin.NameStorageStandalone {
				log.Printf("not found [%s] , use default config [%s]", activePluginName, plugin.NameStorageStandalone)
				activePluginName = plugin.NameStorageStandalone
			}
		} else {
			activePluginName = name
		}
	}

	// 激活插件
	plugin.ActivePlugin(plugin.TypeStorage, activePluginName)
	return nil
}

// 加载连接器插件
func loadConnectorPlugins(cfgVar *gvar.Var) error {
	// 激活插件
	plugin.ActivePlugin(plugin.TypeStorage, plugin.NameConnectorStandalone)
	return nil
}

// 注册插件
func registerPlugins() {
	controller.RegisterStoragePlugin()
	admin.RegisterStoragePlugin()

	// controller.RegisterConnectorPlugin()
}

// InitPliguins 初始化插件
func InitPliguins() error {
	for typeName, plugins := range plugin.GetPlugins() {
		log.Printf("init %s plugins", typeName)
		for name, p := range plugins {
			log.Printf("init %s plugin", name)
			err := p.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
