package model

import (
	"time"
)

var tableName = "topic"

// Topic 主题
type Topic struct {
	Id         string    `bson:"id" json:"id"`                   // ID
	Name       string    `bson:"name" json:"name"`               // 名称
	CreateTime time.Time `bson:"create_time" json:"create_time"` // 创建时间
}

// topicColumns 主题表所有列信息
type topicColumns struct {
	Id         string // ID
	Name       string // 名称
	CreateTime string // 创建时间
}

// topicInfo 主题信息
type topicInfo struct {
	table   string
	columns topicColumns
}

var TopicInfo = topicInfo{
	table: tableName,
	columns: topicColumns{
		Id:         "id",
		Name:       "name",
		CreateTime: "create_time",
	},
}

// Table 表名
func (info *topicInfo) Table() string {
	return info.table
}

// Columns 字段名
func (info *topicInfo) Columns() topicColumns {
	return info.columns
}
