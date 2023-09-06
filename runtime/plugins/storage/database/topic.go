package database

import (
	"context"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"time"
)

type topicService struct{}

var tService = new(topicService)

// QueryByName 根据名称查询
func (s *topicService) QueryByName(ctx context.Context, name string) (topic *model.Topic, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = DB(ctx, model.TopicInfo.Table()).Where(model.TopicInfo.Columns().Name, name).Scan(&topic)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Create 创建主题
func (s *topicService) Create(ctx context.Context, topic *model.Topic) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.TopicInfo.Table()).Insert(topic)
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
		dao := DB(ctx, model.TopicInfo.Table())
		if name != "" {
			dao = dao.WhereLike(model.TopicInfo.Columns().Name, "%"+name+"%")
		}
		cnt, err := dao.Count()
		if err != nil {
			g.Throw(err)
		}
		count = int64(cnt)
		if offset >= 0 && limit > 0 {
			dao = dao.Offset(offset).Limit(limit)
		}
		err = dao.OrderDesc(model.TopicInfo.Columns().CreateTime).Scan(&topics)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// DeleteById 根据ID删除主题
func (s *topicService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.TopicInfo.Table()).Where(model.TopicInfo.Columns().Id, id).Delete()
		if err != nil {
			g.Throw(err)
		}
	})
	return
}
