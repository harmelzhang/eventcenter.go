package mongodb

import (
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/storage"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gctx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type plugin struct{}

func init() {
	plugins.Register(plugins.NameStorageMongodb, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeStorage
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) error {
	clientOptions := options.Client().ApplyURI(config["uri"].String())
	ctx := gctx.New()
	conn, err := mongo.Connect(ctx, clientOptions)
	poolSize := config["poolSize"].Uint64()
	if poolSize > 0 {
		clientOptions.SetMaxPoolSize(poolSize)
	}
	if err != nil {
		return err
	}
	err = conn.Ping(ctx, nil)
	if err != nil {
		return err
	}
	InitDB(conn.Database(config["database"].String()))
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
