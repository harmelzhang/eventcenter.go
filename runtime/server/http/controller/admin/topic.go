package admin

import (
	"context"
	"eventcenter-go/runtime/model"
	"eventcenter-go/runtime/server/http/api/admin"
	"github.com/google/uuid"
	"time"
)

type topicController struct{}

var TopicController = new(topicController)

// Create 创建主题
func (c topicController) Create(ctx context.Context, req *admin.CreateTopicReq) (resp *admin.CreateTopicRes, err error) {
	topic, err := storagePlugin.TopicService().QueryByName(ctx, req.Name)
	if topic != nil {
		return
	}
	topic = &model.Topic{Id: uuid.NewString(), Name: req.Name, CreateTime: time.Now()}
	err = storagePlugin.TopicService().Create(ctx, topic)
	return
}

// Query 查询主题
func (c topicController) Query(ctx context.Context, req *admin.QueryTopicReq) (resp *admin.QueryTopicRes, err error) {
	resp = new(admin.QueryTopicRes)
	topics, count, err := storagePlugin.TopicService().Query(ctx, req.Name, req.Offset, req.Limit)
	resp.Total = count
	resp.Rows = topics
	return
}

// Delete 删除主题
func (c topicController) Delete(ctx context.Context, req *admin.DeleteTopicReq) (resp *admin.DeleteTopicRes, err error) {
	err = storagePlugin.TopicService().DeleteById(ctx, req.Id)
	return
}
