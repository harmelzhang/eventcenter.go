package rabbitmq

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugins"
	"github.com/gogf/gf/v2/container/gvar"
	"go.uber.org/atomic"
)

type plugin struct {
	started atomic.Bool
}

func init() {
	plugins.Register(plugins.NameConnectorRabbitMQ, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeConnector
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) error {
	return nil
}

// Producer 获取生产者
func (p *plugin) Producer() (connector.Producer, error) {
	return nil, nil
}

// Consumer 获取消费者
func (p *plugin) Consumer() (connector.Consumer, error) {
	return nil, nil
}
