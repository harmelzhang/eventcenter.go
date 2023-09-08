package database

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"time"
)

type plugin struct{}

func init() {
	plugins.Register(plugins.NameStorageDB, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeStorage
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) error {
	gdb.SetConfig(gdb.Config{
		plugins.TypeStorage: gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:             config["host"].String(),
				Port:             config["port"].String(),
				User:             config["user"].String(),
				Pass:             config["password"].String(),
				Name:             config["name"].String(),
				Type:             config["type"].String(),
				Link:             config["link"].String(),
				Extra:            config["extra"].String(),
				Role:             config["role"].String(),
				Debug:            config["debug"].Bool(),
				Charset:          config["charset"].String(),
				Prefix:           config["prefix"].String(),
				Weight:           config["weight"].Int(),
				MaxIdleConnCount: config["maxIdle"].Int(),
				MaxOpenConnCount: config["maxOpen"].Int(),
				MaxConnLifeTime:  config["maxLifetime"].Duration() * time.Second,
			},
		},
	})
	return nil
}

// TopicService 主题数据访问层
func (p *plugin) TopicService() storage.TopicService {
	return tService
}

// EndpointService 终端数据访问层
func (p *plugin) EndpointService() storage.EndpointService {
	return epService
}

// EventService 事件数据访问层
func (p *plugin) EventService() storage.EventService {
	return eService
}
