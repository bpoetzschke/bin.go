package slack_middleware

import (
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
)

type SlackMiddlewareMock struct {
	mock.Mock
}

func (m *SlackMiddlewareMock) Connect() <-chan *IncomingMessage {
	return m.Called().Get(0).(<-chan *IncomingMessage)
}

func (m *SlackMiddlewareMock) GetBotInfo() *BotInfo {
	return m.Called().Get(0).(*BotInfo)
}

func (m *SlackMiddlewareMock) PostMessage(msg OutgoingMessage) error {
	return m.Called(msg).Error(0)
}

type SlackRTMMock struct {
	mock.Mock
}

func (m *SlackRTMMock) NewConnection(token string, options ...slack.Option) *slack.Client {
	return m.Called(token, options).Get(0).(*slack.Client)
}

func (m *SlackRTMMock) AddReaction(name string, item slack.ItemRef) error {
	return m.Called(name, item).Error(0)
}

func (m *SlackRTMMock) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	called := m.Called(channelID, options)
	return called.String(0), called.String(1), called.Error(2)
}

func (m *SlackRTMMock) ManageConnection() {
	m.Called()
}

func (m *SlackRTMMock) IncomingEvents() chan slack.RTMEvent {
	return m.Called().Get(0).(chan slack.RTMEvent)
}
