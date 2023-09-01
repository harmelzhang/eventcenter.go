package mongodb

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type mongoPlugin struct{}

func init() {
	plugin.Register(plugin.NameStorageMongodb, &mongoPlugin{})
}

func (p *mongoPlugin) Type() string {
	return plugin.TypeStorage
}

func (p *mongoPlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *mongoPlugin) TopicService() storage.TopicService {
	return tService
}

func (p *mongoPlugin) EventService() storage.EventService {
	return nil
}
