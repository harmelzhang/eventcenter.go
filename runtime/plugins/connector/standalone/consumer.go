package standalone

import (
	"errors"
	"eventcenter-go/runtime/connector"
	"fmt"
	"go.uber.org/atomic"
	"log"
	"sync"
)

type consumer struct {
	broker     *Broker
	subscribes map[string]*subscribeWorker
	handler    *connector.EventHandler
	started    atomic.Bool
	mutex      sync.Mutex
}

func NewConsumer() connector.Consumer {
	return &consumer{
		broker:     GetBroker(),
		subscribes: make(map[string]*subscribeWorker),
	}
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
		for topicName := range c.subscribes {
			err = c.Unsubscribe(topicName)
			if err != nil {
				return
			}
			delete(c.subscribes, topicName)
		}
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

	if _, ok := c.subscribes[topicName]; !ok {
		err = c.broker.CreateNewQueueIfAbsent(topicName)
		if err != nil {
			return
		}

		worker := &subscribeWorker{
			broker:    broker,
			topicName: topicName,
			handler:   c.handler,
			quit:      make(chan bool, 1),
		}
		c.subscribes[topicName] = worker

		go worker.run()
	}

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

	if worker, ok := c.subscribes[topicName]; ok {
		delete(c.subscribes, topicName)
		worker.stop()
	}

	return nil
}

// RegisterHandler 注册事件处理器
func (c *consumer) RegisterHandler(handler *connector.EventHandler) {
	c.handler = handler
}

// Worker

type subscribeWorker struct {
	broker    *Broker
	topicName string
	handler   *connector.EventHandler
	quit      chan bool
}

func (worker *subscribeWorker) run() {
	for {
		select {
		case <-worker.quit:
			return
		default:
			err := worker.popMessage()
			if err != nil {
				log.Printf("fail to pop message from broker: %v", err)
				continue
			}
		}
	}
}

func (worker *subscribeWorker) stop() {
	worker.quit <- true
}

func (worker *subscribeWorker) popMessage() (err error) {
	message, err := worker.broker.PopMessage(worker.topicName)

	if err != nil {
		return errors.New(fmt.Sprintf("get message from broker err: %v", err))
	}

	err = worker.handler.Handler(message.event)
	if err != nil {
		return
	}

	return nil
}
