package admin

import (
	"context"
	"errors"
	"eventcenter-go/runtime/server/http/api/admin"
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
	err = storagePlugin.EndpointService().DeleteById(ctx, req.Id)
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

	// TODO 队列操作

	return
}
