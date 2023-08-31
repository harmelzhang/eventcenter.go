package admin

import (
	"context"
	"eventcenter-go/runtime/server/http/api/admin"
)

type eventController struct{}

var EventController = new(eventController)

// Query 查询事件
func (c eventController) Query(ctx context.Context, req *admin.QueryEventReq) (resp *admin.QueryEventRes, err error) {
	return
}

// Delete 删除事件
func (c eventController) Delete(ctx context.Context, req *admin.DeleteEventReq) (resp *admin.DeleteEventRes, err error) {
	err = storagePlugin.EventService().DeleteById(ctx, req.Id)
	return
}
