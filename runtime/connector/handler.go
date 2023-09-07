package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"eventcenter-go/runtime/consts"
	"eventcenter-go/runtime/model"
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/plugins/storage"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
	"net/http"
)

// HandlerFunc 处理函数
type HandlerFunc func(event *cloudevents.Event) (err error)

// EventHandler 事件处理器
type EventHandler struct {
	// Handler 处理函数
	Handler HandlerFunc
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		Handler: func(event *cloudevents.Event) (err error) {
			plugin := plugins.GetActivedPluginByType(plugins.TypeStorage).(storage.Plugin)
			endpointService := plugin.EndpointService()

			ctx := context.TODO()
			endpoints, err := endpointService.QueryByTopicAndType(ctx, event.Subject(), event.Type())
			if err != nil {
				return err
			}

			for _, endpoint := range endpoints {
				switch endpoint.Protocol {
				case consts.ProtocolHTTP:
					go httpHandler(endpoint, event)
				case consts.ProtocolTCP:
					go tcpHandler(endpoint, event)
				case consts.ConfigGrpc:
					go grpcHandler(endpoint, event)
				default:
					log.Printf("not support handler protocol: %s", endpoint.Protocol)
				}
			}

			return nil
		},
	}
}

func httpHandler(endpoint *model.Endpoint, event *cloudevents.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("json marshal err: %v", err)
		return
	}

	_, err = http.Post(endpoint.Endpoint, event.DataContentType(), bytes.NewReader(data))
	if err != nil {
		log.Printf("http handler err: %v", err)
		return
	}
}

func tcpHandler(endpoint *model.Endpoint, event *cloudevents.Event) {
	// TODO ProtocolTCP
}

func grpcHandler(endpoint *model.Endpoint, event *cloudevents.Event) {
	// TODO ProtocolGrpc
}
