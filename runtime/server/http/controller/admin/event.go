package admin

import (
	"context"
	"encoding/json"
	"eventcenter-go/runtime/server/http/api/admin"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"time"
)

type eventController struct{}

var EventController = new(eventController)

// Query 查询事件
func (c eventController) Query(ctx context.Context, req *admin.QueryEventReq) (resp *admin.QueryEventRes, err error) {
	resp = new(admin.QueryEventRes)
	events, count, err := storagePlugin.EventService().Query(ctx, req.Source, req.TopicName, req.Type, req.Offset, req.Limit)
	resp.Total = count
	resp.Rows = events
	return
}

// Delete 删除事件
func (c eventController) Delete(ctx context.Context, req *admin.DeleteEventReq) (resp *admin.DeleteEventRes, err error) {
	err = storagePlugin.EventService().DeleteById(ctx, req.Id)
	return
}

// Create 创建事件
func (c eventController) Create(ctx context.Context, req *admin.CreateEventReq) (resp *admin.CreateEventRes, err error) {
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

	// 入库
	err = eventService.Create(ctx, event)
	if err != nil {
		return
	}

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

	return
}
