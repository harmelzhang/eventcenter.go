package standalone

import (
	"context"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"strings"
	"time"
)

type topicService struct{}

var tService = new(topicService)

// Create 创建主题
func (s topicService) Create(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err = s.QueryByName(ctx, name)
		if err != nil {
			g.Throw(err)
		}
		if topic != nil {
			// g.Throw("存在同名的主题")
			return
		}
		topic = &model.Topic{Id: uuid.NewString(), Name: name, CreateTime: time.Now()}
		cache[topic.Id+":"+topic.Name] = topic
	})
	return
}

// QueryByName 根据名称查询
func (s topicService) QueryByName(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		for key, t := range cache {
			if strings.HasSuffix(key, ":"+name) {
				topic = t
				return
			}
		}
	})
	return
}

// Query 查询主题
func (s topicService) Query(ctx context.Context) (topics []*model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topics = make([]*model.Topic, 0)
		for _, topic := range cache {
			topics = append(topics, topic)
		}
	})
	return
}

// Delete 删除主题
func (s topicService) Delete(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		keys := getKeys()
		for _, key := range keys {
			if strings.HasPrefix(key, id+":") {
				delete(cache, key)
			}
		}
	})
	return
}
