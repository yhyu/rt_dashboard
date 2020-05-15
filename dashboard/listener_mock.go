package dashboard

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDSTopic struct {
	mock.Mock
	EventSchedule time.Duration
}

func (m *MockDSTopic) GetTopic() string {
	args := m.Mock.Called()
	return args.Get(0).(string)
}

func (m *MockDSTopic) Listen(q DSEventQ) error {
	args := m.Mock.Called(q)
	err := args.Error(0)

	c := time.Tick(m.EventSchedule)
	for next := range c {
		q <- next.String()
	}
	return err
}

func (m *MockDSTopic) React(event TopicEvent) (string, error) {
	args := m.Mock.Called(event)
	return args.Get(0).(string), args.Error(1)
}
