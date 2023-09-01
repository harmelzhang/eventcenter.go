package database

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type dataBasePlugin struct{}

func init() {
	plugin.Register(plugin.NameStorageDB, &dataBasePlugin{})
}

func (p *dataBasePlugin) Type() string {
	return plugin.TypeStorage
}

func (p *dataBasePlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *dataBasePlugin) TopicService() storage.TopicService {
	return tService
}

func (p *dataBasePlugin) EventService() storage.EventService {
	return eService
}
