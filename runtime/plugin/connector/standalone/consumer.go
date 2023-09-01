package standalone

import (
	"errors"
	"eventcenter-go/runtime/connector"
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
		broker: GetBroker(),
	}
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
			quit:      make(chan bool),
		}
		c.subscribes[topicName] = worker

		go worker.run()
	}

	return nil
}

func (c *consumer) Unsubscribe(topicName string) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.IsStoped() {
		err = errors.New("consumer is stop unsubscribe topic error")
		return
	}

	// TODO

	return nil
}

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
			err := worker.pollMessage()
			if err != nil {
				log.Printf("fail to poll message from broker, err=%v", err)
				continue
			}
		}
	}
}

func (worker *subscribeWorker) stop() {
	worker.quit <- true
}

func (worker *subscribeWorker) pollMessage() (err error) {
	message, err := worker.broker.GetMessage(worker.topicName)

	if err != nil {
		return errors.New("get message from broker err")
	}

	err = worker.handler.Handler(message.event)
	if err != nil {
		return
	}

	return nil
}
