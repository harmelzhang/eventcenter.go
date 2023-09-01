package storage

import (
	"context"
	"eventcenter-go/runtime/model"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// EventService 事件数据访问层接口
type EventService interface {

	// Create 创建事件
	Create(ctx context.Context, cloudevent cloudevents.Event) (err error)

	// DeleteById 根据ID删除
	DeleteById(ctx context.Context, id string) (err error)

	// Query 查询事件
	Query(ctx context.Context, source, topicName, typ string, offset, limit int) (events []*model.Event, count int64, err error)
}
