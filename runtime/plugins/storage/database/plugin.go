package database

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type dataBasePlugin struct{}

func init() {
	plugins.Register(plugins.NameStorageDB, &dataBasePlugin{})
}

func (p *dataBasePlugin) Type() string {
	return plugins.TypeStorage
}

func (p *dataBasePlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *dataBasePlugin) TopicService() storage.TopicService {
	return tService
}

func (p *dataBasePlugin) EndpointService() storage.EndpointService {
	return epService
}

func (p *dataBasePlugin) EventService() storage.EventService {
	return eService
}
