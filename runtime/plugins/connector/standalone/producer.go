package standalone

import (
	"context"
	"errors"
	"eventcenter-go/runtime/connector"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/atomic"
)

type producer struct {
	broker  *Broker
	started atomic.Bool
}

func NewProducer() connector.Producer {
	return &producer{
		broker: GetBroker(),
	}
}

func (p *producer) Init() error {
	return nil
}

func (p *producer) IsStarted() bool {
	return p.started.Load()
}

func (p *producer) IsStoped() bool {
	return !p.started.Load()
}

func (p *producer) Start() error {
	p.started.CAS(false, true)
	return nil
}

func (p *producer) Stop() error {
	p.started.CAS(true, false)
	return nil
}

func (p *producer) Publish(ctx context.Context, event *cloudevents.Event) (err error) {
	if p.IsStoped() {
		err = errors.New("producer is stop publish event error")
		return
	}

	_, err = p.broker.PutMessage(event.Subject(), event)
	if err != nil {
		return
	}

	return nil
}
