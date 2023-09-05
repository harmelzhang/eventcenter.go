package mongodb

import (
	"context"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type topicService struct{}

var tService = new(topicService)

// QueryByName 根据名称查询
func (s *topicService) QueryByName(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = DB(ctx, model.TopicInfo.Table()).QuerySet().Q(model.TopicInfo.Columns().Name, name).One(&topic)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Create 创建主题
func (s *topicService) Create(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic = &model.Topic{Id: uuid.NewString(), Name: name, CreateTime: time.Now()}
		_, err = DB(ctx, model.TopicInfo.Table()).InsertOne(ctx, topic)
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
		qs := DB(ctx, model.TopicInfo.Table()).QuerySet()
		if name != "" {
			qs.Filter(bson.D{
				{model.TopicInfo.Columns().Name, primitive.Regex{Pattern: name, Options: "i"}},
			})
		}
		count, err = qs.Count()
		if err != nil {
			g.Throw(err)
		}
		if offset >= 0 && limit > 0 {
			qs.Skip(int64(offset)).Limit(int64(limit))
		}
		err = qs.Sort(bson.E{Key: model.TopicInfo.Columns().CreateTime, Value: -1}).All(&topics)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// DeleteById 根据ID删除主题
func (s *topicService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.TopicInfo.Table()).DeleteOne(ctx, bson.M{"id": id})
		if err != nil {
			g.Throw(err)
		}
	})
	return
}
