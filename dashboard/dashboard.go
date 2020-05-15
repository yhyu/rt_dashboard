package dashboard

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

const (
	defaultMaxBacklog        = 10
	defaultPublishLatestData = true
)

var (
	dsConns      = map[string]map[DSClientId]*websocket.Conn{}
	connLock     sync.Mutex
	lastClientId = DSClientId(0)

	latestData = map[string][]byte{}

	maxBacklog        = defaultMaxBacklog
	publishLatestData = defaultPublishLatestData
)

// DSClientId dashboard client id
type DSClientId uint32

// DSOptions dashboard configuration options
type DSOptions struct {
	DataBacklog       *int
	PublishLatestData *bool
}

// SetOptions sets dashboard configuration options
func SetOptions(options *DSOptions) {
	if options == nil {
		return
	}

	if options.DataBacklog != nil {
		maxBacklog = *options.DataBacklog
	}

	if options.PublishLatestData != nil {
		publishLatestData = *options.PublishLatestData
	}
}

// NewDashboard
func NewDashboard(topic string, conn *websocket.Conn) error {
	id, err := registerConnection(topic, conn)
	if err != nil {
		return err
	}

	// publish latest data
	if publishLatestData {
		data := latestData[topic]
		if data != nil && len(data) > 0 {
			conn.Write(data)
		}
	}

	return startClient(topic, id, conn)
}

// registerConnection registers a client connection
func registerConnection(topic string, conn *websocket.Conn) (DSClientId, error) {
	if _, ok := dsTopics[topic]; !ok {
		return DSClientId(0), fmt.Errorf("invalid topic '%s'", topic)
	}

	connLock.Lock()
	defer connLock.Unlock()
	topicConns, ok := dsConns[topic]
	if !ok {
		topicConns = map[DSClientId]*websocket.Conn{}
		dsConns[topic] = topicConns
	}

	lastClientId++
	topicConns[lastClientId] = conn

	log.Printf("topic '%s' add client '%d', has #of clients: %d\n", topic, lastClientId, len(topicConns))
	return lastClientId, nil
}

func unregisterConnection(topic string, id DSClientId) error {
	connLock.Lock()
	defer connLock.Unlock()
	topicConns, ok := dsConns[topic]
	if !ok {
		return fmt.Errorf("invalid topic '%s'", topic)
	}

	if _, ok := topicConns[id]; !ok {
		return fmt.Errorf("client '%d' not found", id)
	}

	delete(topicConns, id)
	log.Printf("topic '%s' remove client '%d', remains #of clients: %d\n", topic, id, len(topicConns))
	return nil
}

func startClient(topic string, id DSClientId, conn *websocket.Conn) error {
	// wait for client colse
	for {
		msg := make([]byte, 512)
		if _, err := conn.Read(msg); err != nil {
			log.Println("[Closed] read real-time status:", err)
			break
		}
		// received something else: just ignore
	}

	// delete client
	if err := unregisterConnection(topic, id); err != nil {
		return fmt.Errorf("unregister client '%d' from topic '%s' fail: %v\n",
			id, topic, err)
	}
	return nil
}

// pushDataToClients pushes data to clients
func pushDataToClients(topic string, data []byte) error {
	// catch latest data
	if publishLatestData {
		latestData[topic] = data
	}

	connLock.Lock()
	defer connLock.Unlock()
	topicConns, ok := dsConns[topic]
	if !ok {
		// nobody is interested in the topic
		return nil
	}

	for _, q := range topicConns {
		// async to avoid client no response
		go func(conn *websocket.Conn) {
			conn.Write(data)
		}(q)
	}
	return nil
}
