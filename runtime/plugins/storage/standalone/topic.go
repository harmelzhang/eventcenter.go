package standalone

import (
	"context"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"sort"
	"strings"
	"time"
)

type topicService struct{}

var tService = new(topicService)

// QueryByName 根据名称查询
func (s *topicService) QueryByName(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		for key, t := range topicCache {
			if strings.HasSuffix(key, ":"+name) {
				topic = t
				return
			}
		}
	})
	return
}

// Create 创建主题
func (s *topicService) Create(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic = &model.Topic{Id: uuid.NewString(), Name: name, CreateTime: time.Now()}
		topicCache[topic.Id+":"+topic.Name] = topic
	})
	return
}

// QueryOrCreateByName 根据名称查询，如果查询不到则创建
func (s *topicService) QueryOrCreateByName(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err = s.QueryByName(ctx, name)
		if err != nil {
			g.Throw(err)
		}
		if topic == nil {
			topic, err = s.Create(ctx, name)
			if err != nil {
				g.Throw(err)
			}
		}
	})
	return
}

// Query 查询主题
func (s *topicService) Query(ctx context.Context, name string, offset, limit int) (topics []*model.Topic, count int64, err error) {
	topics = make([]*model.Topic, 0)
	err = g.Try(ctx, func(ctx context.Context) {
		for _, topic := range topicCache {
			topics = append(topics, topic)
		}

		if name != "" {
			filter := make([]*model.Topic, 0)
			for _, topic := range topics {
				if strings.Contains(topic.Name, name) {
					filter = append(filter, topic)
				}
			}
			topics = filter
		}

		// 倒序排序
		sort.Slice(topics, func(i, j int) bool {
			return topics[i].CreateTime.Unix() > topics[j].CreateTime.Unix()
		})

		if offset >= 0 && limit > 0 {
			if offset >= len(topics) {
				topics = make([]*model.Topic, 0)
			} else {
				end := offset + limit
				if end > len(topics) {
					end = len(topics)
				}
				topics = topics[offset:end]
			}
		}
	})
	return
}

// DeleteById 根据ID删除主题
func (s *topicService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		keys := getKeys(typeTopic)
		for _, key := range keys {
			if strings.HasPrefix(key, id+":") {
				delete(topicCache, key)
			}
		}
	})
	return
}
