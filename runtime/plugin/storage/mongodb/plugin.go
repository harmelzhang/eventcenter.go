package mongodb

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
)

type mongoPlugin struct{}

func init() {
	plugin.Register(plugin.NameStorageMongodb, &mongoPlugin{})
}

func (p *mongoPlugin) Type() string {
	return plugin.TypeStorage
}

func (p *mongoPlugin) Init() error {
	return nil
}

func (p *mongoPlugin) TopicService() storage.TopicService {
	return tService
}

func (p *mongoPlugin) EventService() storage.EventService {
	return nil
}
