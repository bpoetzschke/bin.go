package slack_middleware

import (
	"github.com/nlopes/slack"
)

type BotInfo struct {
	Name string
	ID   string
}

type MessageType string

const (
	MessageTypeDirectMessage  = MessageType("direct_message")
	MessageTypeSelfMessage    = MessageType("self_message")
	MessageTypeChannelMessage = MessageType("channel_message")
	MessageTypeUnknowMessage  = MessageType("unknown")
)

type Message struct {
	Type      MessageType
	Message   string
	Timestamp string
	Channel   string
	rtm       *slack.RTM
}
