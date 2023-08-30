package standalone

import "eventcenter-go/runtime/model"

var cache = make(map[string]*model.Topic)

func getKeys() (keys []string) {
	for key, _ := range cache {
		keys = append(keys, key)
	}
	return keys
}
