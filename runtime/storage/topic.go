package storage

import (
	"context"
	"eventcenter-go/runtime/model"
)

// TopicService 主题数据访问层接口
type TopicService interface {

	// QueryByName 根据名称查询
	QueryByName(ctx context.Context, name string) (topic *model.Topic, err error)

	// Create 创建主题
	Create(ctx context.Context, name string) (topic *model.Topic, err error)

	// QueryOrCreateByName 根据名称查询，如果查询不到则创建
	QueryOrCreateByName(ctx context.Context, name string) (topic *model.Topic, err error)

	// Query 查询主题
	Query(ctx context.Context, name string, offset, limit int) (topics []*model.Topic, count int64, err error)

	// DeleteById 根据ID删除主题
	DeleteById(ctx context.Context, id string) (err error)
}
