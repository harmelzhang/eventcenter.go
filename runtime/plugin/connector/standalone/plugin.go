package standalone

import (
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugin"
	"github.com/gogf/gf/v2/container/gvar"
)

type standalonePlugin struct {
	consumer connector.Consumer
	producer connector.Producer
}

func init() {
	plugin.Register(plugin.NameConnectorStandalone, &standalonePlugin{})
}

func (p *standalonePlugin) Type() string {
	return plugin.TypeConnector
}

func (p *standalonePlugin) Init(config map[string]*gvar.Var) (err error) {
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

func (p *standalonePlugin) Producer() (connector.Producer, error) {
	return p.producer, nil
}

func (p *standalonePlugin) Consumer() (connector.Consumer, error) {
	return p.consumer, nil
}
