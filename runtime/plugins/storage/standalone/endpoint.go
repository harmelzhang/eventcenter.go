package standalone

import (
	"context"
	"eventcenter-go/runtime/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"sort"
	"strings"
	"time"
)

type endpointService struct{}

var epService = new(endpointService)

// Create 创建终端
func (s *endpointService) Create(ctx context.Context, serverName, topicName, typ, protocol, endpoint string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		topic, err := tService.QueryOrCreateByName(ctx, topicName)
		if err != nil {
			g.Throw(err)
		}

		endpoint := &model.Endpoint{
			Id:           uuid.NewString(),
			ServerName:   serverName,
			TopicId:      topic.Id,
			Type:         typ,
			Protocol:     protocol,
			Endpoint:     endpoint,
			RegisterTime: time.Now(),
		}

		endpointCache[endpoint.Id] = endpoint
	})
	return
}

// DeleteById 根据ID删除
func (s *endpointService) DeleteById(ctx context.Context, id string) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		keys := getKeys(typeEndpoint)
		for _, key := range keys {
			if key == id {
				delete(endpointCache, key)
			}
		}
	})
	return
}

// Update 更新终端
func (s *endpointService) Update(ctx context.Context, endpoint *model.Endpoint) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		endpointCache[endpoint.Id] = endpoint
	})
	return
}

// Query 查询终端
func (s *endpointService) Query(ctx context.Context, serverName, topicName, typ, protocol string, offset, limit int) (endpoints []*model.Endpoint, count int64, err error) {
	endpoints = make([]*model.Endpoint, 0)
	err = g.Try(ctx, func(ctx context.Context) {
		for _, endpoint := range endpointCache {
			endpoints = append(endpoints, endpoint)
		}

		if serverName != "" {
			filter := make([]*model.Endpoint, 0)
			for _, endpoint := range endpoints {
				if strings.Contains(endpoint.ServerName, serverName) {
					filter = append(filter, endpoint)
				}
			}
			endpoints = filter
		}
		if topicName != "" {
			topics, _, err := tService.Query(ctx, topicName, 0, -1)
			if err != nil {
				g.Throw(err)
			}
			topicIds := make([]string, 0)
			for _, topic := range topics {
				topicIds = append(topicIds, topic.Id)
			}
			// 过滤
			filter := make([]*model.Endpoint, 0)
			for _, endpoint := range endpoints {
				for _, topicId := range topicIds {
					if endpoint.TopicId == topicId {
						filter = append(filter, endpoint)
						break
					}
				}
			}
			endpoints = filter
		}
		if typ != "" {
			filter := make([]*model.Endpoint, 0)
			for _, endpoint := range endpoints {
				if strings.Contains(endpoint.Type, typ) {
					filter = append(filter, endpoint)
				}
			}
			endpoints = filter
		}
		if protocol != "" {
			filter := make([]*model.Endpoint, 0)
			for _, endpoint := range endpoints {
				if strings.Contains(endpoint.Protocol, protocol) {
					filter = append(filter, endpoint)
				}
			}
			endpoints = filter
		}

		count = int64(len(endpoints))

		// 倒序排序
		sort.Slice(endpoints, func(i, j int) bool {
			return endpoints[i].RegisterTime.Unix() > endpoints[j].RegisterTime.Unix()
		})

		if offset >= 0 && limit > 0 {
			if offset >= len(endpoints) {
				endpoints = make([]*model.Endpoint, 0)
			} else {
				end := offset + limit
				if end > len(endpoints) {
					end = len(endpoints)
				}
				endpoints = endpoints[offset:end]
			}
		}
	})
	return
}

// QueryByTopicAndServer 根据主题和服务查询
func (s *endpointService) QueryByTopicAndServer(ctx context.Context, topicName, typ, serverName, protocol string) (endpoint *model.Endpoint, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		endpoints := make([]*model.Endpoint, 0)
		for _, ep := range endpointCache {
			endpoints = append(endpoints, ep)
		}
		if len(endpoints) == 0 {
			return
		}

		topic, err := tService.QueryOrCreateByName(ctx, topicName)
		if err != nil {
			g.Throw(err)
		}

		for _, ep := range endpoints {
			if ep.TopicId == topic.Id && ep.Type == typ && ep.ServerName == serverName && ep.Protocol == protocol {
				endpoint = ep
				return
			}
		}
	})
	return
}
