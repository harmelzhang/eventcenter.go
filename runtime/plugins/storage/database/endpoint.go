package database

import (
	"context"
	"eventcenter-go/runtime/model"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

type endpointService struct{}

var epService = new(endpointService)

// Create 创建终端
func (s *endpointService) Create(ctx context.Context, endpoint *model.Endpoint) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.EndpointInfo.Table()).Insert(endpoint)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// DeleteById 根据ID删除
func (s *endpointService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.EndpointInfo.Table()).Where(model.EndpointInfo.Columns().Id, id).Delete()
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Update 更新终端
func (s *endpointService) Update(ctx context.Context, endpoint *model.Endpoint) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.EndpointInfo.Table()).Where(model.EndpointInfo.Columns().Id, endpoint.Id).Update(endpoint)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Query 查询终端
func (s *endpointService) Query(ctx context.Context, serverName, topicName, typ, protocol string, offset, limit int) (endpoints []*model.Endpoint, count int64, err error) {
	endpoints = make([]*model.Endpoint, 0)
	err = g.Try(ctx, func(ctx context.Context) {
		dao := DB(ctx, model.EndpointInfo.Table())

		if serverName != "" {
			dao = dao.WhereLike(model.EndpointInfo.Columns().ServerName, "%"+serverName+"%")
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
			dao = dao.Where(model.EndpointInfo.Columns().TopicId+" in (?)", topicIds)
		}
		if typ != "" {
			dao = dao.WhereLike(model.EndpointInfo.Columns().Type, "%"+typ+"%")
		}
		if protocol != "" {
			dao = dao.WhereLike(model.EndpointInfo.Columns().Protocol, "%"+protocol+"%")
		}

		cnt, err := dao.Count()
		if err != nil {
			g.Throw(err)
		}
		count = int64(cnt)

		if offset >= 0 && limit > 0 {
			dao = dao.Offset(offset).Limit(limit)
		}

		err = dao.OrderDesc(model.TopicInfo.Columns().CreateTime).Scan(&endpoints)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// QueryById 根据ID查询
func (s *endpointService) QueryById(ctx context.Context, id string) (endpoint *model.Endpoint, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = DB(ctx, model.EndpointInfo.Table()).Where(model.EndpointInfo.Columns().Id, id).Scan(&endpoint)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// QueryByTopicAndServer 根据主题和服务查询
func (s *endpointService) QueryByTopicAndServer(ctx context.Context, topicName, typ, serverName, protocol string) (endpoint *model.Endpoint, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err := tService.QueryOrCreateByName(ctx, topicName)
		if err != nil {
			g.Throw(err)
		}

		dao := DB(ctx, model.EndpointInfo.Table())
		dao = dao.Where(model.EndpointInfo.Columns().TopicId, topic.Id)
		dao = dao.Where(model.EndpointInfo.Columns().Type, typ)
		dao = dao.Where(model.EndpointInfo.Columns().ServerName, serverName)
		dao = dao.Where(model.EndpointInfo.Columns().Protocol, protocol)

		err = dao.Scan(&endpoint)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// QueryByTopicAndType 根据主题和类型查询
func (s *endpointService) QueryByTopicAndType(ctx context.Context, topicName, typ string) (endpoints []*model.Endpoint, err error) {
	endpoints = make([]*model.Endpoint, 0)
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err := tService.QueryByName(ctx, topicName)
		if err != nil {
			g.Throw(err)
		}
		if topic == nil {
			g.Throw(fmt.Sprintf("not found topic [%s]", topicName))
		}

		dao := DB(ctx, model.EndpointInfo.Table())
		dao = dao.Where(model.EndpointInfo.Columns().TopicId, topic.Id)
		dao = dao.Where(model.EndpointInfo.Columns().Type, typ)

		err = dao.Scan(&endpoints)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// QueryCountByTopic 根据主题查询数量
func (s *endpointService) QueryCountByTopic(ctx context.Context, topicName string) (count int64, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		dao := DB(ctx, model.EndpointInfo.Table())
		if topicName != "" {
			topics, _, err := tService.Query(ctx, topicName, 0, -1)
			if err != nil {
				g.Throw(err)
			}
			topicIds := make([]string, 0)
			for _, topic := range topics {
				topicIds = append(topicIds, topic.Id)
			}
			dao = dao.Where(model.EndpointInfo.Columns().TopicId+" in (?)", topicIds)
		}
		cnt, err := dao.Count()
		if err != nil {
			g.Throw(err)
		}
		count = int64(cnt)
	})
	return
}
