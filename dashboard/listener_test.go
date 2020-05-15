package dashboard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ListenerTestSuite struct {
	suite.Suite
}

func TestListenerTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) SetupSuite() {}

func (s *ListenerTestSuite) Test_RegisterTopic() {
	mockConn := &MockDSTopic{
		EventSchedule: 500 * time.Millisecond,
	}
	mockConn.On("GetTopic").Return("my_topic")
	mockConn.On("React", mock.Anything).Return("results", nil)
	mockConn.On("Listen", mock.Anything).Return(nil)
	err := RegisterTopic(mockConn)
	assert.NoError(s.T(), err)

	err = RegisterTopic(mockConn)
	assert.Error(s.T(), err)

	time.Sleep(time.Second)
}
