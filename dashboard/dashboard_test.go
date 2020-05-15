package dashboard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	ws "golang.org/x/net/websocket"
)

type DashboardTestSuite struct {
	suite.Suite
}

// TestGeneralTestSuite invokes a normal testcase for general adapter testsuite.
func TestDashboardTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardTestSuite))
}

func (s *DashboardTestSuite) SetupSuite() {}

func (s *DashboardTestSuite) Test_SetOptions() {
	SetOptions(nil)
	assert.Equal(s.T(), maxBacklog, defaultMaxBacklog)
	assert.Equal(s.T(), publishLatestData, defaultPublishLatestData)

	newMaxBacklog := defaultMaxBacklog * 2
	newPublishLatestData := !defaultPublishLatestData
	options := &DSOptions{
		DataBacklog:       &newMaxBacklog,
		PublishLatestData: &newPublishLatestData,
	}
	SetOptions(options)
	assert.Equal(s.T(), maxBacklog, *options.DataBacklog)
	assert.Equal(s.T(), publishLatestData, *options.PublishLatestData)
}

func (s *DashboardTestSuite) Test_001_NewDashboard() {
	err := NewDashboard("ds_topic", &ws.Conn{})
	assert.Error(s.T(), err)
}

func (s *DashboardTestSuite) Test_002_registerConnection() {
	_, err := registerConnection("reg_topic", &ws.Conn{})
	assert.Error(s.T(), err)

	mockConn := &MockDSTopic{
		EventSchedule: time.Hour,
	}
	mockConn.On("GetTopic").Return("reg_topic")
	mockConn.On("React", mock.Anything).Return("results", nil)
	mockConn.On("Listen", mock.Anything).Return(nil)
	err = RegisterTopic(mockConn)
	assert.NoError(s.T(), err)

	id, err := registerConnection("reg_topic", &ws.Conn{})
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), uint32(id), uint32(0))
}

func (s *DashboardTestSuite) Test_003_unregisterConnection() {
	err := unregisterConnection("unknown", DSClientId(123))
	assert.Error(s.T(), err)

	err = unregisterConnection("reg_topic", DSClientId(123))
	assert.Error(s.T(), err)

	id, err := registerConnection("reg_topic", &ws.Conn{})
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), uint32(id), uint32(0))

	err = unregisterConnection("reg_topic", id)
	assert.NoError(s.T(), err)
}
