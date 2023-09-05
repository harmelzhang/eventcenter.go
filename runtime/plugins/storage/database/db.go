package database

import (
	"context"
	"eventcenter-go/runtime/plugins"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func DB(ctx context.Context, table string) *gdb.Model {
	return g.DB(plugins.TypeStorage).Model(table).Safe().Ctx(ctx)
}
