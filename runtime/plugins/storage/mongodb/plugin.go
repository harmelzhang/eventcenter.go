package mongodb

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type mongoPlugin struct{}

func init() {
	plugins.Register(plugins.NameStorageMongodb, &mongoPlugin{})
}

func (p *mongoPlugin) Type() string {
	return plugins.TypeStorage
}

func (p *mongoPlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *mongoPlugin) TopicService() storage.TopicService {
	return tService
}

func (p *mongoPlugin) EndpointService() storage.EndpointService {
	return epService
}

func (p *mongoPlugin) EventService() storage.EventService {
	return nil
}
