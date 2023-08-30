package mongodb

import (
	"context"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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
		_, err = DB(ctx, model.TopicInfo.Table()).InsertOne(ctx, topic)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// QueryByName 根据名称查询
func (s topicService) QueryByName(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = DB(ctx, model.TopicInfo.Table()).QuerySet().Q(model.TopicInfo.Columns().Name, name).One(&topic)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Query 查询主题
func (s topicService) Query(ctx context.Context) (topics []*model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = DB(ctx, model.TopicInfo.Table()).QuerySet().Sort(bson.E{Key: model.TopicInfo.Columns().CreateTime, Value: -1}).All(&topics)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Delete 删除主题
func (s topicService) Delete(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.TopicInfo.Table()).DeleteOne(ctx, bson.M{"id": id})
		if err != nil {
			g.Throw(err)
		}
	})
	return
}
