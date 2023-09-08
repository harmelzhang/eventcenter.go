package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"eventcenter-go/runtime/connector"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/streadway/amqp"
	"go.uber.org/atomic"
)

type producer struct {
	config  map[string]*gvar.Var
	started atomic.Bool
}

func NewProducer(config map[string]*gvar.Var) connector.Producer {
	return &producer{config: config}
}

// Init 初始化
func (p *producer) Init() error {
	conn, err := amqp.Dial(p.config["uri"].String())
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	err = ch.ExchangeDeclare(
		p.config["exchange"].String(),
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

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

	conn, err := amqp.Dial(p.config["uri"].String())
	if err != nil {
		return
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		return
	}
	defer func() { _ = ch.Close() }()

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	err = ch.Publish(p.config["exchange"].String(), event.Subject(), false, false, amqp.Publishing{
		Body: data,
	})

	return err
}
