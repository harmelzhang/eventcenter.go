package mongodb

import (
	"context"
	"encoding/json"
	"eventcenter-go/runtime/model"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/frame/g"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		_, err = DB(ctx, model.EventInfo.Table()).InsertOne(ctx, event)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// DeleteById 根据ID删除
func (s *eventService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.EventInfo.Table()).DeleteOne(ctx, bson.M{"id": id})
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
		qs := DB(ctx, model.EventInfo.Table()).QuerySet()

		if source != "" {
			qs.Filter(bson.D{
				{model.EventInfo.Columns().Source, primitive.Regex{Pattern: source, Options: "i"}},
			})
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
			qs.Q(model.EventInfo.Columns().TopicId, bson.M{"$in": topicIds})
		}
		if typ != "" {
			qs.Filter(bson.D{
				{model.EventInfo.Columns().Type, primitive.Regex{Pattern: typ, Options: "i"}},
			})
		}

		count, err = qs.Count()
		if err != nil {
			g.Throw(err)
		}

		if offset >= 0 && limit > 0 {
			qs.Skip(int64(offset)).Limit(int64(limit))
		}

		err = qs.Sort(bson.E{Key: model.EventInfo.Columns().CreateTime, Value: -1}).All(&events)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}
