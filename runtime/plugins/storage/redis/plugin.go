package redis

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"time"
)

type plugin struct{}

func init() {
	plugins.Register(plugins.NameStorageRedis, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeStorage
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) error {
	gredis.SetConfig(&gredis.Config{
		Address:         config["address"].String(),
		Pass:            config["password"].String(),
		Db:              config["db"].Int(),
		User:            config["user"].String(),
		MinIdle:         config["minIdle"].Int(),
		MaxIdle:         config["maxIdle"].Int(),
		MaxActive:       config["maxActive"].Int(),
		MaxConnLifetime: config["maxConnLifetime"].Duration() * time.Second,
		IdleTimeout:     config["idleTimeout"].Duration() * time.Second,
		WaitTimeout:     config["waitTimeout"].Duration() * time.Second,
		DialTimeout:     config["dialTimeout"].Duration() * time.Second,
		ReadTimeout:     config["readTimeout"].Duration() * time.Second,
		WriteTimeout:    config["writeTimeout"].Duration() * time.Second,
		MasterName:      config["masterName"].String(),
		TLS:             config["tls"].Bool(),
		TLSSkipVerify:   config["tlsSkipVerify"].Bool(),
		SlaveOnly:       config["slaveOnly"].Bool(),
	}, plugins.TypeStorage)
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
