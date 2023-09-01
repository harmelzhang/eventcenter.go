package standalone

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
var evventKeyPrefix = "event"

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

		eventCache[event.Id] = event
	})
	return
}

// DeleteById 根据ID删除
func (s *eventService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		keys := getKeys(typeEvent)
		for _, key := range keys {
			if key == id {
				delete(eventCache, key)
			}
		}
	})
	return
}

// Query 查询事件
func (s *eventService) Query(ctx context.Context, source, topicName, typ string, offset, limit int) (events []*model.Event, count int64, err error) {
	events = make([]*model.Event, 0)
	err = g.Try(ctx, func(ctx context.Context) {
		for _, event := range eventCache {
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
