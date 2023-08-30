package standalone

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
)

type standalonePlugin struct{}

func init() {
	plugin.Register(plugin.NameStorageStandalone, &standalonePlugin{})
}

func (p *standalonePlugin) Type() string {
	return plugin.TypeStorage
}

func (p *standalonePlugin) Init() error {
	return nil
}

func (p *standalonePlugin) TopicService() storage.TopicService {
	return tService
}
