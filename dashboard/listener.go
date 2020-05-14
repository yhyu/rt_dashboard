package dashboard

import (
	"fmt"
	"log"
)

var (
	dsTopics = map[string]DSTopic{}
)

// TopicEvent topic change event data type
type TopicEvent interface{}

// DSEventQ topic event queue
type DSEventQ chan TopicEvent

// DSTopic defines dashboard topic interface
type DSTopic interface {
	// GetTopic gets topic id
	GetTopic() string

	// Listen listens for event change, any change should be pushed to DSEventQ
	Listen(q DSEventQ) error

	// React reacts for event change, and returns processed data to push to frontend
	React(event TopicEvent) (string, error)
}

// RegisterTopic registers a topic
func RegisterTopic(topic DSTopic) error {
	topicName := topic.GetTopic()
	if _, ok := dsTopics[topicName]; ok {
		return fmt.Errorf("duplicated topic '%s'", topicName)
	}

	dsTopics[topicName] = topic

	// start listen
	go topicListener(topic)
	return nil
}

func topicListener(topic DSTopic) {
	eventQ := make(DSEventQ, maxBacklog)

	go func() {
		for event := range eventQ {
			data, err := topic.React(event)
			if err != nil {
				log.Printf("topic '%s' react event fail: %v\n", topic.GetTopic(), err)
				continue
			}

			if err := pushDataToClients(topic.GetTopic(), []byte(data)); err != nil {
				log.Printf("topic '%s' push data to client fail: %v\n", topic.GetTopic(), err)
				continue
			}
		}
	}()

	if err := topic.Listen(eventQ); err != nil {
		log.Printf("topic '%s' listen fail: %v\n", topic.GetTopic(), err)
	}
}
