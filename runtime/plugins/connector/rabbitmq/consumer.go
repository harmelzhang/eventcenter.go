package rabbitmq

import (
	"errors"
	"eventcenter-go/runtime/connector"
	"go.uber.org/atomic"
	"sync"
)

type consumer struct {
	handler *connector.EventHandler
	started atomic.Bool
	mutex   sync.Mutex
}

func NewConsumer() connector.Consumer {
	return &consumer{}
}

func (c *consumer) Init() error {
	return nil
}

func (c *consumer) IsStarted() bool {
	return c.started.Load()
}

func (c *consumer) IsStoped() bool {
	return !c.started.Load()
}

func (c *consumer) Start() error {
	c.started.CAS(false, true)
	return nil
}

func (c *consumer) Stop() (err error) {
	if ok := c.started.CAS(true, false); ok {
		// TODO Shutdown
	}

	return nil
}

func (c *consumer) Subscribe(topicName string) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.IsStoped() {
		err = errors.New("consumer is stop subscribe topic error")
		return
	}

	// TODO 订阅

	return nil
}

func (c *consumer) Unsubscribe(topicName string) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.IsStoped() {
		err = errors.New("consumer is stop unsubscribe topic error")
		return
	}

	// TODO 取消订阅

	return nil
}

func (c *consumer) RegisterHandler(handler *connector.EventHandler) {
	c.handler = handler
}
