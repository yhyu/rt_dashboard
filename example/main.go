package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	ds "../dashboard"

	"golang.org/x/net/websocket"
)

var (
	srvUrl = flag.String("u", ":54321", "server url")
)

const (
	myTopic = "topic1"
)

func main() {
	flag.Parse()

	if err := ds.RegisterTopic(&MyListener{time.Now()}); err != nil {
		log.Fatalln("register my topic fail:", err)
	}

	// websocket service
	http.Handle("/ws/topic1", websocket.Handler(DashboardServiceTopic1))
	log.Fatal(http.ListenAndServe(*srvUrl, nil))
}

func DashboardServiceTopic1(conn *websocket.Conn) {
	if err := ds.NewDashboard(myTopic, conn); err != nil {
		log.Println("create new topic1 dashboard fail:", err)
	}
}

type MyListener struct {
	StartTime time.Time
}

func (l *MyListener) GetTopic() string {
	return myTopic
}

func (l *MyListener) Listen(q ds.DSEventQ) error {
	// waiting for something...
	// simulate event generation: generate int64 data every 30 seconds
	c := time.Tick(30 * time.Second)
	for next := range c {
		q <- next.Unix()
	}
	return nil
}

func (l *MyListener) React(event ds.TopicEvent) (string, error) {
	// do something to react the event
	// simulate action: calculate time elapsed time in seconds
	data, ok := event.(int64)
	if !ok {
		return "", fmt.Errorf("invalid event: %v", event)
	}
	return strconv.FormatInt(data-l.StartTime.Unix(), 10), nil
}
