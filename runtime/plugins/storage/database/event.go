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
func (s *eventService) Create(ctx context.Context, cloudevent cloudevents.Event) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err := tService.QueryOrCreateByName(ctx, cloudevent.Subject())
		if err != nil {
			g.Throw(err)
		}

		bytes, err := json.Marshal(cloudevent)
		if err != nil {
			g.Throw(err)
		}

		event := &model.Event{
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
func (s *eventService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.EventInfo.Table()).Where(model.EventInfo.Columns().Id, id).Delete()
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Query 查询事件
func (s *eventService) Query(ctx context.Context, source, topicName, typ string, offset, limit int) (events []*model.Event, count int64, err error) {
	events = make([]*model.Event, 0)
	err = g.Try(ctx, func(ctx context.Context) {
		dao := DB(ctx, model.EventInfo.Table())

		if source != "" {
			dao = dao.WhereLike(model.EventInfo.Columns().Source, "%"+source+"%")
		}
		if topicName != "" {
			topics, _, err := tService.Query(ctx, topicName, 0, -1)
			if err != nil {
				g.Throw(err)
			}
			topicIds := make([]string, 0)
			for _, topic := range topics {
				topicIds = append(topicIds, topic.Id)
			}
			dao = dao.Where(model.EventInfo.Columns().TopicId+" in (?)", topicIds)
		}
		if typ != "" {
			dao = dao.WhereLike(model.EventInfo.Columns().Type, "%"+typ+"%")
		}

		cnt, err := dao.Count()
		if err != nil {
			g.Throw(err)
		}
		count = int64(cnt)

		if offset >= 0 && limit > 0 {
			dao = dao.Offset(offset).Limit(limit)
		}

		err = dao.OrderDesc(model.TopicInfo.Columns().CreateTime).Scan(&events)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}
