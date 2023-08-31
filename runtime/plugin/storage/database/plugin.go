package database

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
)

type dataBasePlugin struct{}

func init() {
	plugin.Register(plugin.NameStorageDB, &dataBasePlugin{})
}

func (p *dataBasePlugin) Type() string {
	return plugin.TypeStorage
}

func (p *dataBasePlugin) Init() error {
	return nil
}

func (p *dataBasePlugin) TopicService() storage.TopicService {
	return tService
}

func (p *dataBasePlugin) EventService() storage.EventService {
	return eService
}
