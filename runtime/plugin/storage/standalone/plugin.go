package standalone

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type standalonePlugin struct{}

func init() {
	plugin.Register(plugin.NameStorageStandalone, &standalonePlugin{})
}

func (p *standalonePlugin) Type() string {
	return plugin.TypeStorage
}

func (p *standalonePlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *standalonePlugin) TopicService() storage.TopicService {
	return tService
}

func (p *standalonePlugin) EventService() storage.EventService {
	return eService
}
