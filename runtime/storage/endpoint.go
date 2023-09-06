package storage

import (
	"context"
	"eventcenter-go/runtime/model"
)

// EndpointService 终端数据访问层接口
type EndpointService interface {

	// Create 创建终端
	Create(ctx context.Context, serverName, topicName, typ, protocol, endpoint string) (ep *model.Endpoint, err error)

	// DeleteById 根据ID删除
	DeleteById(ctx context.Context, id string) (err error)

	// Update 更新终端
	Update(ctx context.Context, endpoint *model.Endpoint) (err error)

	// Query 查询终端
	Query(ctx context.Context, serverName, topicName, typ, protocol string, offset, limit int) (endpoints []*model.Endpoint, count int64, err error)

	// QueryByTopicAndServer 根据主题和服务查询
	QueryByTopicAndServer(ctx context.Context, topicName, typ, serverName, protocol string) (endpoint *model.Endpoint, err error)
}
