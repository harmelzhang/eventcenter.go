package storage

import (
	"context"
	"eventcenter-go/runtime/model"
)

// TopicService 主题数据访问层接口
type TopicService interface {
	// Create 创建主题
	Create(ctx context.Context, name string) (topic *model.Topic, err error)

	// QueryByName 根据名称查询
	QueryByName(ctx context.Context, name string) (topic *model.Topic, err error)

	// Query 查询主题
	Query(ctx context.Context) (topics []*model.Topic, err error)

	// Delete 删除主题
	Delete(ctx context.Context, id string) (err error)
}
