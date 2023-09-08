package redis

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugins"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"time"
)

type plugin struct {
	consumer connector.Consumer
	producer connector.Producer
}

func init() {
	plugins.Register(plugins.NameConnectorRedis, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeConnector
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) (err error) {
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
	}, plugins.TypeConnector)

	p.consumer = NewConsumer()
	err = p.consumer.Start()
	if err != nil {
		return
	}

	p.producer = NewProducer()
	err = p.producer.Start()
	if err != nil {
		return
	}

	return nil
}

// Producer 获取生产者
func (p *plugin) Producer() (connector.Producer, error) {
	return p.producer, nil
}

// Consumer 获取消费者
func (p *plugin) Consumer() (connector.Consumer, error) {
	return p.consumer, nil
}
