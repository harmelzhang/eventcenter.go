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

type endpointService struct{}

var epService = new(endpointService)

// Create 创建终端
func (s *endpointService) Create(ctx context.Context, serverName, topicName, protocol, endpoint string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err := tService.QueryOrCreateByName(ctx, topicName)
		if err != nil {
			g.Throw(err)
		}

		endpoint := &model.Endpoint{
			Id:           uuid.NewString(),
			ServerName:   serverName,
			TopicId:      topic.Id,
			Protocol:     protocol,
			Endpoint:     endpoint,
			RegisterTime: time.Now(),
		}
		_, err = DB(ctx, model.EndpointInfo.Table()).InsertOne(ctx, endpoint)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// DeleteById 根据ID删除
func (s *endpointService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = DB(ctx, model.EndpointInfo.Table()).DeleteOne(ctx, bson.M{"id": id})
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Update 更新终端
func (s *endpointService) Update(ctx context.Context, endpoint *model.Endpoint) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		filter := bson.M{"id": endpoint.Id}
		doc := bson.M{
			"$set": bson.M{
				"server_name": endpoint.ServerName,
				"topic_id":    endpoint.TopicId,
				"protocol":    endpoint.Protocol,
				"endpoint":    endpoint.Endpoint,
			},
		}
		_, err = DB(ctx, model.EndpointInfo.Table()).UpdateOne(ctx, filter, doc)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// Query 查询终端
func (s *endpointService) Query(ctx context.Context, serverName, topicName, protocol string, offset, limit int) (endpoints []*model.Endpoint, count int64, err error) {
	endpoints = make([]*model.Endpoint, 0)
	err = g.Try(ctx, func(ctx context.Context) {
		qs := DB(ctx, model.EndpointInfo.Table()).QuerySet()

		if serverName != "" {
			qs.Filter(bson.D{
				{model.EndpointInfo.Columns().ServerName, primitive.Regex{Pattern: serverName, Options: "i"}},
			})
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
			qs.Q(model.EndpointInfo.Columns().TopicId, bson.M{"$in": topicIds})
		}
		if protocol != "" {
			qs.Filter(bson.D{
				{model.EndpointInfo.Columns().Protocol, primitive.Regex{Pattern: protocol, Options: "i"}},
			})
		}

		count, err = qs.Count()
		if err != nil {
			g.Throw(err)
		}

		if offset >= 0 && limit > 0 {
			qs.Skip(int64(offset)).Limit(int64(limit))
		}

		err = qs.Sort(bson.E{Key: model.EndpointInfo.Columns().RegisterTime, Value: -1}).All(&endpoints)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}

// QueryByTopicAndServer 根据主题和服务查询
func (s *endpointService) QueryByTopicAndServer(ctx context.Context, topicName, serverName, protocol string) (endpoint *model.Endpoint, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err := tService.QueryOrCreateByName(ctx, topicName)
		if err != nil {
			g.Throw(err)
		}

		qs := DB(ctx, model.EndpointInfo.Table()).QuerySet()
		qs.Q(model.EndpointInfo.Columns().TopicId, topic.Id)
		qs.Q(model.EndpointInfo.Columns().ServerName, serverName)
		qs.Q(model.EndpointInfo.Columns().Protocol, protocol)

		err = qs.One(&endpoint)
		if err != nil {
			g.Throw(err)
		}
	})
	return
}
