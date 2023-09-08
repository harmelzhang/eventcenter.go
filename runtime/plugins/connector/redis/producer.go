package redis

import (
	"context"
	"errors"
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugins"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/frame/g"
	"go.uber.org/atomic"
)

type producer struct {
	//broker  *Broker
	started atomic.Bool
}

func NewProducer() connector.Producer {
	return &producer{
		//broker: GetBroker(),
	}
}

// Init 初始化
func (p *producer) Init() error {
	return nil
}

// IsStarted 是否启动
func (p *producer) IsStarted() bool {
	return p.started.Load()
}

// IsStoped 是否停止
func (p *producer) IsStoped() bool {
	return !p.started.Load()
}

// Start 启动服务
func (p *producer) Start() error {
	p.started.CAS(false, true)
	return nil
}

// Stop 停止服务
func (p *producer) Stop() error {
	p.started.CAS(true, false)
	return nil
}

// Publish 发布事件
func (p *producer) Publish(ctx context.Context, event *cloudevents.Event) (err error) {
	if p.IsStoped() {
		err = errors.New("producer is stop publish event error")
		return
	}

	_, err = g.Redis(plugins.TypeConnector).Publish(ctx, event.Subject(), event)
	if err != nil {
		return
	}

	return nil
}
