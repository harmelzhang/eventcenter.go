package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseRes 基础返回
type BaseRes struct {
	g.Meta `mime:"application/json"`
}

// PageRes 分页返回
type PageRes struct {
	BaseRes
	Total int         `json:"total"`
	Rows  interface{} `json:"rows"`
}
