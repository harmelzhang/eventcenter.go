package standalone

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"sync"
)

// Message 消息
type Message struct {
	event *cloudevents.Event
}

// MessageQueue 消息队列
type MessageQueue struct {
	items  []*Message
	mutex  sync.Mutex
	newMsg sync.Cond
}

func NewMessageQueue() (*MessageQueue, error) {
	queue := &MessageQueue{
		items: make([]*Message, 0),
	}
	queue.newMsg = sync.Cond{L: &queue.mutex}

	return queue, nil
}

func (queue *MessageQueue) Put(message *Message) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.items = append(queue.items, message)
	queue.newMsg.Signal()
}

func (queue *MessageQueue) Pop() *Message {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	for len(queue.items) == 0 {
		queue.newMsg.Wait()
	}

	message := queue.items[0]
	queue.items = queue.items[1:]

	return message
}

// Broker 消息代理
type Broker struct {
	queueContainer map[string]*MessageQueue // 队列容器：主题名:队列
}

var once sync.Once
var broker *Broker

func GetBroker() *Broker {
	once.Do(func() {
		broker = &Broker{
			queueContainer: make(map[string]*MessageQueue),
		}
	},
	)
	return broker
}

func (b *Broker) CreateNewQueueIfAbsent(topicName string) (err error) {
	if _, ok := b.queueContainer[topicName]; ok {
		return
	}

	queue, err := NewMessageQueue()
	if err != nil {
		return
	}

	b.queueContainer[topicName] = queue

	return
}

func (b *Broker) PutMessage(topicName string, event *cloudevents.Event) (message *Message, err error) {
	if err = b.CreateNewQueueIfAbsent(topicName); err != nil {
		return
	}

	message = &Message{event: event}
	b.queueContainer[topicName].Put(message)

	return
}

func (b *Broker) PopMessage(topicName string) (*Message, error) {
	if err := b.CreateNewQueueIfAbsent(topicName); err != nil {
		return nil, err
	}
	return b.queueContainer[topicName].Pop(), nil
}
