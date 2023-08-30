package redis

import (
	"context"
	"encoding/json"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"strings"
	"time"
)

type topicService struct{}

var tService = new(topicService)

// topic:[ID]:[NAME] -> META DATA
var topicKeyPrefix = "topic"

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

// QueryByName 根据名称查询
func (s topicService) QueryByName(ctx context.Context, name string) (topic *model.Topic, err error) {
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
		return
	})
	return
}

// Query 查询主题
func (s topicService) Query(ctx context.Context) (topics []*model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		keys, err := DB(ctx).Keys(ctx, topicKeyPrefix+":*")
		if err != nil {
			g.Throw(err)
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
	})
	return
}

// Delete 删除主题
func (s topicService) Delete(ctx context.Context, id string) (err error) {
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
