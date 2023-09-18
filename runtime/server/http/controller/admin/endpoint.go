package admin

import (
	"context"
	"errors"
	"eventcenter-go/runtime/model"
	"eventcenter-go/runtime/server/http/api/admin"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
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
	// 校验地址
	if req.IsMicro == 1 {
		if !strings.HasPrefix(req.Url, "/") {
			err = errors.New("事件处理地址格式错误，必须为绝对路径")
			return
		}
	} else {
		if req.Protocol != "http" {
			err = errors.New("事件处理协议暂只支持HTTP")
			return
		}
		if !(strings.HasPrefix(req.Url, "http://") || strings.HasPrefix(req.Url, "https://")) {
			err = errors.New("事件处理地址协议暂只支持HTTP和HTTPS")
			return
		}
	}

	topicService := storagePlugin.TopicService()
	endpointService := storagePlugin.EndpointService()

	topic, err := topicService.QueryOrCreateByName(ctx, req.TopicName)
	if err != nil {
		return
	}

	endpoint, err := endpointService.QueryByTopicAndServer(ctx, topic.Name, req.Type, req.ServerName, req.Protocol)
	if err != nil {
		return
	}

	if endpoint != nil {
		err = errors.New("该服务下存在同主题同类型的处理终端")
		return
	}

	endpoint = &model.Endpoint{
		Id:           uuid.NewString(),
		ServerName:   req.ServerName,
		IsMicro:      req.IsMicro,
		TopicId:      topic.Id,
		Type:         req.Type,
		Protocol:     req.Protocol,
		Endpoint:     req.Url,
		RegisterTime: time.Now(),
	}
	err = endpointService.Create(ctx, endpoint)
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
