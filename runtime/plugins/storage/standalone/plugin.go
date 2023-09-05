package standalone

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type standalonePlugin struct{}

func init() {
	plugins.Register(plugins.NameStorageStandalone, &standalonePlugin{})
}

func (p *standalonePlugin) Type() string {
	return plugins.TypeStorage
}

func (p *standalonePlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *standalonePlugin) TopicService() storage.TopicService {
	return tService
}

func (p *standalonePlugin) EndpointService() storage.EndpointService {
	return epService
}

func (p *standalonePlugin) EventService() storage.EventService {
	return eService
}
