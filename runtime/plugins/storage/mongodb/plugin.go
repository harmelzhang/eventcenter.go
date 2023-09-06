package mongodb

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type plugin struct{}

func init() {
	plugins.Register(plugins.NameStorageMongodb, &plugin{})
}

func (p *plugin) Type() string {
	return plugins.TypeStorage
}

func (p *plugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *plugin) TopicService() storage.TopicService {
	return tService
}

func (p *plugin) EndpointService() storage.EndpointService {
	return epService
}

func (p *plugin) EventService() storage.EventService {
	return nil
}
