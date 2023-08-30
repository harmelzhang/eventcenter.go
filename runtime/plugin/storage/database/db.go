package database

import (
	"context"
	"eventcenter-go/runtime/plugin"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func DB(ctx context.Context, table string) *gdb.Model {
	return g.DB(plugin.TypeStorage).Model(table).Safe().Ctx(ctx)
}
