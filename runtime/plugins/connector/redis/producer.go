package redis

import (
	"context"
	"encoding/json"
	"errors"
	"eventcenter-go/runtime/connector"
	"eventcenter-go/runtime/plugins"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/frame/g"
	"go.uber.org/atomic"
)

type producer struct {
	started     atomic.Bool
	queuePrefix string
}

func NewProducer(queuePrefix string) connector.Producer {
	return &producer{queuePrefix: queuePrefix}
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
func (p *producer) Publish(ctx context.Context, cloudevent *cloudevents.Event) (err error) {
	if p.IsStoped() {
		err = errors.New("producer is stop publish event error")
		return
	}

	bytes, err := json.Marshal(cloudevent)
	if err != nil {
		return
	}

	key := fmt.Sprintf("%s:%s", p.queuePrefix, cloudevent.Subject())
	_, err = g.Redis(plugins.TypeConnector).LPush(ctx, key, string(bytes))
	if err != nil {
		return
	}

	return nil
}
