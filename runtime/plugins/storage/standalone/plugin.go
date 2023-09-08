package standalone

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
)

type plugin struct{}

func init() {
	plugins.Register(plugins.NameStandalone, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeStorage
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) error {
	return nil
}

// TopicService 主题数据访问层
func (p *plugin) TopicService() storage.TopicService {
	return tService
}

// EndpointService 终端数据访问层
func (p *plugin) EndpointService() storage.EndpointService {
	return epService
}

// EventService 事件数据访问层
func (p *plugin) EventService() storage.EventService {
	return eService
}
