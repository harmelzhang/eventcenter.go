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

// Init 初始化
func (c *consumer) Init() error {
	return nil
}

// IsStarted 是否启动
func (c *consumer) IsStarted() bool {
	return c.started.Load()
}

// IsStoped 是否停止
func (c *consumer) IsStoped() bool {
	return !c.started.Load()
}

// Start 启动服务
func (c *consumer) Start() error {
	c.started.CAS(false, true)
	return nil
}

// Stop 停止服务
func (c *consumer) Stop() (err error) {
	if ok := c.started.CAS(true, false); ok {
		// TODO Shutdown
	}

	return nil
}

// Subscribe 订阅
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

// Unsubscribe 取消订阅
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

// RegisterHandler 注册事件处理器
func (c *consumer) RegisterHandler(handler *connector.EventHandler) {
	c.handler = handler
}
