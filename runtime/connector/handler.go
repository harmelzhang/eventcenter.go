package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"eventcenter-go/runtime/consts"
	"eventcenter-go/runtime/model"
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/plugins/registry"
	"eventcenter-go/runtime/plugins/storage"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"io"
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
				case consts.ProtocolHTTP, consts.ProtocolHTTPS:
					go httpHandler(endpoint, event)
				case consts.ProtocolTCP:
					go tcpHandler(endpoint, event)
				case consts.ProtocolGrpc:
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

	url := endpoint.Endpoint
	if endpoint.IsMicro == 1 {
		registryPlugin := plugins.GetActivedPluginByType(plugins.TypeRegistry)
		if registryPlugin == nil {
			log.Printf("not found actived registry plugin")
			return
		}
		registryService := registryPlugin.(registry.Plugin).Service()
		ins, err := registryService.FindService(endpoint.ServerName)
		if err != nil {
			log.Printf("registry find service err: %v", err)
			return
		}
		if ins == nil {
			log.Println("registry not found service")
			return
		}
		if ins.Port == 80 {
			url = fmt.Sprintf("%s://%s%s", endpoint.Protocol, ins.Address, endpoint.Endpoint)
		} else {
			url = fmt.Sprintf("%s://%s:%d%s", endpoint.Protocol, ins.Address, ins.Port, endpoint.Endpoint)
		}
	}
	httpResp, err := http.Post(url, event.DataContentType(), bytes.NewReader(data))
	if err != nil {
		log.Printf("http handler err: %v", err)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("read http response body err: %v", err)
		return
	}
	log.Printf("http response: %d -> %s", httpResp.StatusCode, string(body))
}

func tcpHandler(endpoint *model.Endpoint, event *cloudevents.Event) {
	// TODO ProtocolTCP
}

func grpcHandler(endpoint *model.Endpoint, event *cloudevents.Event) {
	// TODO ProtocolGrpc
}
