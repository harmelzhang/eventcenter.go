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

func (p *plugin) Type() string {
	return plugins.TypeConnector
}

func (p *plugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *plugin) Producer() (connector.Producer, error) {
	return nil, nil
}

func (p *plugin) Consumer() (connector.Consumer, error) {
	return nil, nil
}
