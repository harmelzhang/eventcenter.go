package database

import (
	"context"
	"encoding/json"
	"eventcenter-go/runtime/model"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/frame/g"
)

type eventService struct{}

var eService = new(eventService)

// Create 创建事件
func (s eventService) Create(ctx context.Context, cloudevent cloudevents.Event) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err := tService.QueryOrCreateByName(ctx, cloudevent.Subject())
		if err != nil {
			g.Throw(err)
		}

		bytes, err := json.Marshal(cloudevent)
		if err != nil {
			g.Throw(err)
		}

		event := model.Event{
			Id:          cloudevent.ID(),
			Source:      cloudevent.Source(),
			TopicId:     topic.Id,
			Type:        cloudevent.Type(),
			Data:        string(cloudevent.Data()),
			CreateTime:  cloudevent.Time(),
			CloudEvents: string(bytes),
		}
		_, err = DB(ctx, model.EventInfo.Table()).Insert(event)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// DeleteById 根据ID删除
func (s eventService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.EventInfo.Table()).Where(model.EventInfo.Columns().Id, id).Delete()
		if err != nil {
			g.Throw(err)
		}
	})
	return
}
