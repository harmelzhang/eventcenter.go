package controller

import (
	"context"
	"errors"
	"eventcenter-go/runtime/server/http/api"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
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

	// 如果没有该主题则创建
	topic, err := storagePlugin.TopicService().QueryByName(ctx, req.TopicName)
	if err != nil {
		return
	}
	if topic == nil {
		topic, err = storagePlugin.TopicService().Create(ctx, req.TopicName)
		if err != nil {
			return
		}
	}

	// TODO 注册入库

	return
}

// Trigger 触发
func (c processController) Trigger(ctx context.Context, req *api.TriggerReq) (resp *api.TriggerRes, err error) {
	uid := uuid.NewString()

	event := cloudevents.NewEvent()
	event.SetSubject(req.TopicName)
	event.SetTime(time.Now())
	event.SetID(uid)

	// TODO 丢入队列

	resp = new(api.TriggerRes)
	resp.EventId = uid
	return
}
