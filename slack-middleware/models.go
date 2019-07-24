package slack_middleware

import (
	"github.com/nlopes/slack"
)

type SlackRTM interface {
	NewConnection(token string, options ...slack.Option) *slack.Client
	AddReaction(name string, item slack.ItemRef) error
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	ManageConnection()
	IncomingEvents() chan slack.RTMEvent
}

type BotInfo struct {
	Name string
	ID   string
}

type MessageType string

const (
	MessageTypeDirectMessage  = MessageType("direct_message")
	MessageTypeSelfMessage    = MessageType("self_message")
	MessageTypeChannelMessage = MessageType("channel_message")
	MessageTypeUnknownMessage = MessageType("unknown")
)

type BaseMessage struct {
	Message string
	Channel string
}

type IncomingMessage struct {
	BaseMessage
	Type      MessageType
	Timestamp string
	UserID    string
	rtm       SlackRTM
}

type OutgoingMessage struct {
	BaseMessage
	Attachments []string
}
