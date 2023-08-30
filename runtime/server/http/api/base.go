package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// EmptyRes 空返回
type EmptyRes struct {
	g.Meta `mime:"application/json"`
}

// PageRes 分页返回
type PageRes struct {
	EmptyRes
	Total int         `json:"total"`
	Rows  interface{} `json:"rows"`
}
