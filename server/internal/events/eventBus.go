package events

import (
	"log"
	"sync"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

type DataChannel chan DataEvent

type DataChannelSlice []DataChannel

type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()

	if chans, found := eb.subscribers[topic]; found {

		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

func (eb *EventBus) Unsubscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if chans, found := eb.subscribers[topic]; found {
		for i, c := range chans {
			if c == ch {
				eb.subscribers[topic] = append(chans[:i], chans[i+1:]...)
				break
			}
		}
		if len(eb.subscribers[topic]) == 0 {
			delete(eb.subscribers, topic)
		}
		log.Println("unsubscribed from channel :::", &ch)
	}
	eb.rm.Unlock()
}

var EB = &EventBus{
	subscribers: map[string]DataChannelSlice{},
}
