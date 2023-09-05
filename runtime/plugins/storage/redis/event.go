package redis

import (
	"context"
	"encoding/json"
	"eventcenter-go/runtime/model"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/frame/g"
	"sort"
	"strings"
)

type eventService struct{}

var eService = new(eventService)

// event:[ID] -> DATA
var eventKeyPrefix = "event"

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
		value, err := json.Marshal(event)
		if err != nil {
			g.Throw(err)
		}

		_, err = DB(ctx).Set(ctx, eventKeyPrefix+":"+event.Id, string(value))
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// DeleteById 根据ID删除
func (s *eventService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx).Del(ctx, eventKeyPrefix+":"+id)
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
		keys, err := DB(ctx).Keys(ctx, eventKeyPrefix+":*")
		if err != nil {
			g.Throw(err)
		}
		if len(keys) == 0 {
			return
		}

		values, err := DB(ctx).MGet(ctx, keys...)
		if err != nil {
			g.Throw(err)
		}

		for _, value := range values {
			event := new(model.Event)
			err = json.Unmarshal(value.Bytes(), event)
			if err != nil {
				g.Throw(err)
			}
			events = append(events, event)
		}

		if source != "" {
			filter := make([]*model.Event, 0)
			for _, event := range events {
				if strings.Contains(event.Source, source) {
					filter = append(filter, event)
				}
			}
			events = filter
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
			// 过滤
			filter := make([]*model.Event, 0)
			for _, event := range events {
				for _, topicId := range topicIds {
					if event.TopicId == topicId {
						filter = append(filter, event)
						break
					}
				}
			}
			events = filter
		}
		if typ != "" {
			filter := make([]*model.Event, 0)
			for _, event := range events {
				if strings.Contains(event.Type, typ) {
					filter = append(filter, event)
				}
			}
			events = filter
		}

		count = int64(len(events))

		// 倒序排序
		sort.Slice(events, func(i, j int) bool {
			return events[i].CreateTime.Unix() > events[j].CreateTime.Unix()
		})

		if offset >= 0 && limit > 0 {
			if offset >= len(events) {
				events = make([]*model.Event, 0)
			} else {
				end := offset + limit
				if end > len(events) {
					end = len(events)
				}
				events = events[offset:end]
			}
		}
	})
	return
}
