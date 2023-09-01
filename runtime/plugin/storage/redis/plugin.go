package redis

import (
	"eventcenter-go/runtime/plugin"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type redisPlugin struct{}

func init() {
	plugin.Register(plugin.NameStorageRedis, &redisPlugin{})
}

func (p *redisPlugin) Type() string {
	return plugin.TypeStorage
}

func (p *redisPlugin) Init(config map[string]*gvar.Var) error {
	return nil
}

func (p *redisPlugin) TopicService() storage.TopicService {
	return tService
}

func (p *redisPlugin) EventService() storage.EventService {
	return eService
}
