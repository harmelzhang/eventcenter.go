package standalone

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugin"
	"github.com/gogf/gf/v2/container/gvar"
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

func (p *standalonePlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *standalonePlugin) Producer() (connector.Producer, error) {
	return nil, nil
}

func (p *standalonePlugin) Consumer() (connector.Consumer, error) {
	return nil, nil
}
