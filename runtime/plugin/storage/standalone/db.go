package standalone

import "eventcenter-go/runtime/model"

type modelType string

const (
	typeTopic modelType = "topic"
	typeEvent           = "event"
)

var topicCache = make(map[string]*model.Topic)
var eventCache = make(map[string]*model.Event)

// 获取所有 Key
func getKeys(typ modelType) (keys []string) {
	if typ == typeTopic {
		for key, _ := range topicCache {
			keys = append(keys, key)
		}
		return keys
	} else if typ == typeEvent {
		for key, _ := range eventCache {
			keys = append(keys, key)
		}
		return keys
	}
	return keys
}
