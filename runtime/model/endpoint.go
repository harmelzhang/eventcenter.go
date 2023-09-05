package model

import "time"

var endpointTableName = "endpoint"

// Endpoint 终端
type Endpoint struct {
	Id           string    `bson:"id" json:"id"`                       // ID
	ServerName   string    `bson:"server_name" json:"server_name"`     // 服务名
	TopicId      string    `bson:"topic_id" json:"topic_id"`           // 主题ID
	Protocol     string    `bson:"protocol" json:"protocol"`           // 处理协议
	Endpoint     string    `bson:"endpoint" json:"endpoint"`           // 终端地址
	RegisterTime time.Time `bson:"register_time" json:"register_time"` // 注册时间
}

// endpointColumns 终端表所有列信息
type endpointColumns struct {
	Id           string // ID
	ServerName   string // 服务名
	TopicId      string // 主题ID
	Protocol     string // 协议
	Endpoint     string // 终端地址
	RegisterTime string // 注册时间
}

// endpointInfo 终端表信息
type endpointInfo struct {
	table   string
	columns endpointColumns
}

var EndpointInfo = endpointInfo{
	table: endpointTableName,
	columns: endpointColumns{
		Id:           "id",
		ServerName:   "server_name",
		TopicId:      "topic_id",
		Protocol:     "protocol",
		Endpoint:     "endpoint",
		RegisterTime: "register_time",
	},
}

// Table 数据表名
func (info *endpointInfo) Table() string {
	return info.table
}

// Columns 字段名
func (info *endpointInfo) Columns() endpointColumns {
	return info.columns
}
