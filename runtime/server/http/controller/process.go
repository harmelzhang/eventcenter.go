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

	endpoint, err := endpointService.QueryByTopicAndServer(ctx, topic.Name, req.ServerName, req.Protocol)
	if err != nil {
		return
	}

	if endpoint == nil {
		// 创建
		err = endpointService.Create(ctx, req.ServerName, topic.Name, req.Protocol, req.Url)
		if err != nil {
			return
		}
		// TODO 队列操作
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

	endpoint, err := endpointService.QueryByTopicAndServer(ctx, req.TopicName, req.ServerName, req.Protocol)
	if err != nil {
		return
	}

	err = endpointService.DeleteById(ctx, endpoint.Id)

	// TODO 队列操作

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

	// TODO 队列操作

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
