# Realtime Dashboard
A library to create realtime dashboard through websocket.

## Prerequisite
This library uses weksocket package.

```
go get golang.org/x/net/websocket
```

## How to use?

1. Implement DSTopic interface to define a dashboard topic behavior: 

```
// DSTopic defines dashboard topic interface
type DSTopic interface {
	// GetTopic gets topic id
	GetTopic() string

	// Listen listens for event change, any change should be pushed to DSEventQ
	Listen(q DSEventQ) error

	// React reacts for event change, and returns processed data taht frontend will see
	React(event TopicEvent) (string, error)
}
```

2. Implement websocket handler
```
func DashboardServiceTopic1(conn *websocket.Conn) {
	if err := ds.NewDashboard(myTopic, conn); err != nil {
		log.Println("create topic1 dashboard fail:", err)
	}
}
```

3. Register topic and websocket handler
```
	RegisterTopic(&MyDSTopic{})
	http.Handle("/ws/topic1", websocket.Handler(DashboardServiceTopic1))
	log.Fatal(http.ListenAndServe("127.0.0.1:54321", nil))
```
Notes: Please see the [example](https://github.com/yhyu/rt_dashaboard/tree/master/example)
