package standalone

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugin"
	"go.uber.org/atomic"
)

type standalonePlugin struct {
	started atomic.Bool
}

func init() {
	plugin.Register(plugin.NameConnectorStandalone, &standalonePlugin{})
}

func (p *standalonePlugin) Type() string {
	return plugin.TypeConnector
}

func (p *standalonePlugin) Init() error {
	return nil
}

func (p *standalonePlugin) Producer() (connector.Producer, error) {
	return nil, nil
}

func (p *standalonePlugin) Consumer() (connector.Consumer, error) {
	return nil, nil
}
