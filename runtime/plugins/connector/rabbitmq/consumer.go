package rabbitmq

import (
	"encoding/json"
	"errors"
	"eventcenter-go/runtime/connector"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/streadway/amqp"
	"go.uber.org/atomic"
	"log"
	"sync"
)

type consumer struct {
	config     map[string]*gvar.Var
	subscribes map[string]*subscribeWorker
	handler    *connector.EventHandler
	started    atomic.Bool
	mutex      sync.Mutex
}

func NewConsumer(config map[string]*gvar.Var) connector.Consumer {
	return &consumer{config: config, subscribes: make(map[string]*subscribeWorker)}
}

// Init 初始化
func (c *consumer) Init() error {
	conn, err := amqp.Dial(c.config["uri"].String())
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
		c.config["exchange"].String(),
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
		worker := &subscribeWorker{
			config:    c.config,
			topicName: topicName,
			handler:   c.handler,
			quit:      make(chan bool, 1),
		}
		c.subscribes[topicName] = worker

		var wg sync.WaitGroup
		wg.Add(1)
		go worker.run(&wg)
		wg.Wait()
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
	config    map[string]*gvar.Var
	topicName string
	handler   *connector.EventHandler
	quit      chan bool
}

func (worker *subscribeWorker) run(wg *sync.WaitGroup) {
	conn, err := amqp.Dial(worker.config["uri"].String())
	if err != nil {
		log.Printf("failed to conn rabbitmq err: %v", err)
		return
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("failed to open a channel: %v", err)
		return
	}
	defer func() { _ = ch.Close() }()

	queue, err := ch.QueueDeclare(
		worker.topicName,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("failed to declare a queue: %s", err)
		return
	}

	err = ch.QueueBind(
		queue.Name,
		worker.topicName,
		worker.config["exchange"].String(),
		false,
		nil,
	)
	if err != nil {
		log.Printf("failed to bind a queue: %s", err)
		return
	}

	delivery, err := ch.Consume(queue.Name, worker.topicName, true, false, false, false, nil)
	if err != nil {
		log.Printf("failed consume messages: %v", err)
		return
	}

	wg.Done()

	for {
		select {
		case <-worker.quit:
			return
		default:
			message, isOK := <-delivery
			if !isOK {
				continue
			}

			if string(message.Body) == "" {
				continue
			}

			event := cloudevents.NewEvent()
			err = json.Unmarshal(message.Body, &event)
			if err != nil {
				log.Printf("fail to unmarshal message err: %v", err)
				continue
			}

			err = worker.handler.Handler(&event)
			if err != nil {
				log.Printf("handler event err: %v", err)
				continue
			}
		}
	}
}

func (worker *subscribeWorker) stop() {
	worker.quit <- true
	// 让最后一个阻塞的消息获取线程处理完
	go func() {
		conn, err := amqp.Dial(worker.config["uri"].String())
		if err != nil {
			log.Printf("failed to conn rabbitmq err: %v", err)
			return
		}
		defer func() { _ = conn.Close() }()

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("failed to open a channel: %v", err)
			return
		}
		defer func() { _ = ch.Close() }()

		err = ch.Publish(worker.config["exchange"].String(), worker.topicName, false, false, amqp.Publishing{
			Body: []byte(""),
		})
		if err != nil {
			log.Printf("publish stop signal message err: %v", err)
		}
	}()
}
