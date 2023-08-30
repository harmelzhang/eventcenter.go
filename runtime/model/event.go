package model

import (
	"time"
)

var eventTableName = "event"

// Event 事件
type Event struct {
	Id          string    `bson:"id" json:"id"`                   // ID
	Source      string    `bson:"source" json:"source"`           // 事件源
	TopicId     string    `bson:"topic_id" json:"topic_id"`       // 主题ID
	Type        string    `bson:"type" json:"type"`               // 事件类型
	Data        string    `bson:"data" json:"data"`               // 上下文数据
	CreateTime  time.Time `bson:"create_time" json:"create_time"` // 创建时间
	CloudEvents string    `bson:"cloudevents" json:"cloudevents"` // 完整CloudEvents规范数据
}

// eventColumns 事件表所有列信息
type eventColumns struct {
	Id          string // ID
	Source      string // 事件源
	TopicId     string // 主题ID
	Type        string // 事件类型
	Data        string // 上下文数据
	CreateTime  string // 创建时间
	CloudEvents string // 完整CloudEvents规范数据
}

// eventInfo 事件表信息
type eventInfo struct {
	table   string
	columns eventColumns
}

var EventInfo = eventInfo{
	table: eventTableName,
	columns: eventColumns{
		Id:          "id",
		Source:      "source",
		TopicId:     "topic_id",
		Type:        "type",
		Data:        "data",
		CreateTime:  "create_time",
		CloudEvents: "cloudevents",
	},
}

// Table 数据表名
func (info *eventInfo) Table() string {
	return info.table
}

// Columns 字段名
func (info *eventInfo) Columns() eventColumns {
	return info.columns
}
