package controller

import (
	"context"
	"encoding/json"
	"errors"
	"eventcenter-go/runtime/server/http/api"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

type processController struct{}

var ProcessController = new(processController)

// Subscribe 订阅
func (c processController) Subscribe(ctx context.Context, req *api.SubscribeReq) (resp *api.SubscribeRes, err error) {
	// 校验地址
	if req.Protocol != "http" {
		err = errors.New("事件处理协议暂只支持HTTP")
		return
	}
	if !(strings.HasPrefix(req.Url, "http://") || strings.HasPrefix(req.Url, "https://")) {
		err = errors.New("事件处理地址协议暂只支持HTTP和HTTPS")
		return
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

	if endpoint == nil {
		// 入库
		_, err = endpointService.Create(ctx, req.ServerName, topic.Name, req.Type, req.Protocol, req.Url)
		if err != nil {
			return
		}
		// 订阅
		consumer, err := connectorPlugin.Consumer()
		if err != nil {
			return resp, err
		}
		err = consumer.Subscribe(topic.Name)
		if err != nil {
			return resp, err
		}
	} else {
		// 更新地址
		if endpoint.Endpoint != req.Url {
			endpoint.Endpoint = req.Url
			err = endpointService.Update(ctx, endpoint)
			if err != nil {
				return
			}
		}
	}

	return
}

// Unsubscribe 取消订阅
func (c processController) Unsubscribe(ctx context.Context, req *api.UnsubscribeReq) (resp *api.UnsubscribeRes, err error) {
	endpointService := storagePlugin.EndpointService()

	endpoint, err := endpointService.QueryByTopicAndServer(ctx, req.TopicName, req.Type, req.ServerName, req.Protocol)
	if err != nil || endpoint == nil {
		return
	}

	err = endpointService.DeleteById(ctx, endpoint.Id)
	if err != nil {
		return
	}

	// 取消订阅
	go func() {
		ctx := context.TODO()
		count, err := endpointService.QueryCountByTopic(ctx, req.TopicName)
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
			err = consumer.Unsubscribe(req.TopicName)
			if err != nil {
				log.Printf("consumer unsubscribe topic [%s] err: %v", req.TopicName, err)
				return
			}
		}
	}()

	return
}

// Trigger 触发
func (c processController) Trigger(ctx context.Context, req *api.TriggerReq) (resp *api.TriggerRes, err error) {
	eventService := storagePlugin.EventService()

	uid := uuid.NewString()

	event := cloudevents.NewEvent()
	event.SetID(uid)
	event.SetSource(req.Source)
	event.SetSubject(req.TopicName)
	event.SetType(req.Type)
	// 数据类型判断
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(req.Data), &data)
	if err != nil {
		err = event.SetData(cloudevents.ApplicationJSON, req.Data)
		if err != nil {
			return
		}
	} else {
		err = event.SetData(cloudevents.ApplicationJSON, data)
		if err != nil {
			return
		}
	}
	event.SetTime(time.Now())

	// 补偿订阅（防止启动时没有对应主题的消费者）
	consumer, err := connectorPlugin.Consumer()
	if err != nil {
		return
	}
	err = consumer.Subscribe(event.Subject())
	if err != nil {
		return
	}

	// 发布
	producer, err := connectorPlugin.Producer()
	if err != nil {
		return
	}
	err = producer.Publish(ctx, &event)
	if err != nil {
		return
	}

	// 入库
	go func() {
		err = eventService.Create(gctx.New(), event)
		if err != nil {
			log.Printf("insert event err: %v", err)
		}
	}()

	resp = new(api.TriggerRes)
	resp.EventId = uid
	return
}
