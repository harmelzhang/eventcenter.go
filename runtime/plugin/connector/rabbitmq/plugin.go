package rabbitmq

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugin"
	"go.uber.org/atomic"
)

type rabbitmqPlugin struct {
	started atomic.Bool
}

func init() {
	plugin.Register(plugin.NameConnectorRabbitMQ, &rabbitmqPlugin{})
}

func (p *rabbitmqPlugin) Type() string {
	return plugin.TypeConnector
}

func (p *rabbitmqPlugin) Init() error {
	return nil
}

func (p *rabbitmqPlugin) Producer() (connector.Producer, error) {
	return nil, nil
}

func (p *rabbitmqPlugin) Consumer() (connector.Consumer, error) {
	return nil, nil
}
