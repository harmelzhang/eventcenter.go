package redis

import (
	"context"
	"encoding/json"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"sort"
	"strings"
	"time"
)

type topicService struct{}

var tService = new(topicService)

// topic:[ID]:[NAME] -> DATA
var topicKeyPrefix = "topic"

// QueryByName 根据名称查询
func (s *topicService) QueryByName(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		keys, err := DB(ctx).Keys(ctx, topicKeyPrefix+":*:"+name)
		if err != nil {
			g.Throw(err)
		}
		for _, key := range keys {
			if strings.HasSuffix(key, ":"+name) {
				value, err := DB(ctx).Get(ctx, key)
				if err != nil {
					g.Throw(err)
				}
				topic = new(model.Topic)
				err = json.Unmarshal(value.Bytes(), topic)
				if err != nil {
					g.Throw(err)
				}
				return
			}
		}
	})
	return
}

// QueryById 根据ID查询
func (s *topicService) QueryById(ctx context.Context, id string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		keys, err := DB(ctx).Keys(ctx, topicKeyPrefix+":"+id+":*")
		if err != nil {
			g.Throw(err)
		}
		for _, key := range keys {
			if strings.HasPrefix(key, ":"+id+":") {
				value, err := DB(ctx).Get(ctx, key)
				if err != nil {
					g.Throw(err)
				}
				topic = new(model.Topic)
				err = json.Unmarshal(value.Bytes(), topic)
				if err != nil {
					g.Throw(err)
				}
				return
			}
		}
	})
	return
}

// Create 创建主题
func (s *topicService) Create(ctx context.Context, topic *model.Topic) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		value, err := json.Marshal(topic)
		if err != nil {
			g.Throw(err)
		}
		_, err = DB(ctx).Set(ctx, topicKeyPrefix+":"+topic.Id+":"+topic.Name, string(value))
		if err != nil {
			g.Throw(err)
		}
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
			topic = &model.Topic{Id: uuid.NewString(), Name: name, CreateTime: time.Now()}
			err = s.Create(ctx, topic)
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
		keys, err := DB(ctx).Keys(ctx, topicKeyPrefix+":*")
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
			topic := new(model.Topic)
			err = json.Unmarshal(value.Bytes(), topic)
			if err != nil {
				g.Throw(err)
			}
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

		count = int64(len(topics))

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
		keys, err := DB(ctx).Keys(ctx, topicKeyPrefix+":"+id+":*")
		if err != nil {
			g.Throw(err)
		}
		for _, key := range keys {
			_, err = DB(ctx).Del(ctx, key)
			if err != nil {
				g.Throw(err)
			}
		}
		return
	})
	return
}
