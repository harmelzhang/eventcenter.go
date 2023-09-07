package admin

import (
	"context"
	"errors"
	"eventcenter-go/runtime/server/http/api/admin"
	"log"
)

type endpointController struct{}

var EndpointController = new(endpointController)

// Query 查询终端
func (c endpointController) Query(ctx context.Context, req *admin.QueryEndpointReq) (resp *admin.QueryEndpointRes, err error) {
	resp = new(admin.QueryEndpointRes)
	endpoints, count, err := storagePlugin.EndpointService().Query(ctx, req.ServerName, req.TopicName, req.Type, req.Protocol, req.Offset, req.Limit)
	resp.Total = count
	resp.Rows = endpoints
	return
}

// Delete 删除终端
func (c endpointController) Delete(ctx context.Context, req *admin.DeleteEndpointReq) (resp *admin.DeleteEndpointRes, err error) {
	endpointService := storagePlugin.EndpointService()
	topicService := storagePlugin.TopicService()

	endpoint, err := endpointService.QueryById(ctx, req.Id)
	if err != nil || endpoint == nil {
		return
	}

	err = endpointService.DeleteById(ctx, req.Id)

	// 取消订阅
	go func() {
		ctx := context.TODO()

		topic, err := topicService.QueryById(ctx, endpoint.TopicId)
		if err != nil {
			log.Printf("query topic [%s] err: %v", endpoint.TopicId, err)
			return
		}

		count, err := endpointService.QueryCountByTopic(ctx, topic.Name)
		if err != nil {
			log.Printf("query count by topic err: %v", err)
			return
		}
		if count == 0 {
			consumer, err := connectorPlugin.Consumer()
			if err != nil {
				log.Printf("connector plugin get consumer err: %v", err)
				return
			}
			err = consumer.Unsubscribe(topic.Name)
			if err != nil {
				log.Printf("consumer unsubscribe topic [%s] err: %v", topic.Name, err)
				return
			}
		}
	}()

	return
}

// Create 创建终端
func (c endpointController) Create(ctx context.Context, req *admin.CreateEndpointReq) (resp *admin.CreateEndpointRes, err error) {
	endpointService := storagePlugin.EndpointService()

	endpoint, err := endpointService.QueryByTopicAndServer(ctx, req.TopicName, req.Type, req.ServerName, req.Protocol)
	if err != nil {
		return
	}

	if endpoint != nil {
		err = errors.New("该服务下存在同主题同类型的处理终端")
		return
	}

	_, err = endpointService.Create(ctx, req.ServerName, req.TopicName, req.Type, req.Protocol, req.Url)
	if err != nil {
		return
	}

	// 订阅
	consumer, err := connectorPlugin.Consumer()
	if err != nil {
		return
	}
	err = consumer.Subscribe(req.TopicName)

	return
}
