package redis

import (
	"context"
	"eventcenter-go/runtime/plugin"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

func DB(ctx context.Context) *gredis.Redis {
	return g.Redis(plugin.TypeStorage)
}
