package rabbitmq

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugins"
	"github.com/gogf/gf/v2/container/gvar"
)

type plugin struct {
	consumer connector.Consumer
	producer connector.Producer
}

func init() {
	plugins.Register(plugins.NameConnectorRabbitMQ, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeConnector
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) (err error) {
	p.consumer = NewConsumer(config)
	err = p.consumer.Init() // 准备 Exchange
	if err != nil {
		return
	}
	err = p.consumer.Start()
	if err != nil {
		return
	}

	p.producer = NewProducer(config)
	err = p.producer.Init() // 准备 Exchange
	if err != nil {
		return
	}
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
