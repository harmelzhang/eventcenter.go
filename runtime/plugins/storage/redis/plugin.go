package redis

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type redisPlugin struct{}

func init() {
	plugins.Register(plugins.NameStorageRedis, &redisPlugin{})
}

func (p *redisPlugin) Type() string {
	return plugins.TypeStorage
}

func (p *redisPlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *redisPlugin) TopicService() storage.TopicService {
	return tService
}

func (p *redisPlugin) EndpointService() storage.EndpointService {
	return epService
}

func (p *redisPlugin) EventService() storage.EventService {
	return eService
}
