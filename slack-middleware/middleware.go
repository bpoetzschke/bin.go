package slack_middleware

import (
	"log"
	"strings"

	"github.com/nlopes/slack"
)

type Middleware interface {
	Connect() <-chan *Message
	GetBotInfo() *BotInfo
}

func NewMiddleware(slackToken string) Middleware {
	mw := middleware{slackToken: slackToken}
	mw.init()

	return &mw
}

type middleware struct {
	slackToken string

	slackApi *slack.Client
	slackRTM *slack.RTM

	eventChannel chan *Message
	botInfo      *BotInfo
	logProvider  *log.Logger
}

func (mw *middleware) init() {
	mw.slackApi = slack.New(mw.slackToken)
	mw.eventChannel = make(chan *Message)
}

func (mw *middleware) Connect() <-chan *Message {
	mw.slackRTM = mw.slackApi.NewRTM()
	go mw.slackRTM.ManageConnection()

	go func() {
		for {
			select {
			case evt := <-mw.slackRTM.IncomingEvents:
				switch evt.Data.(type) {
				case *slack.ConnectedEvent:
					mw.handleConnect(evt.Data)
				case *slack.MessageEvent:
					mw.handleMessageEvent(evt.Data)
				default:
				}

			default:
			}
		}
	}()

	return mw.eventChannel
}

func (mw *middleware) GetBotInfo() *BotInfo {
	return mw.botInfo
}

func (mw *middleware) handleConnect(payload interface{}) {
	connectedEvt := payload.(*slack.ConnectedEvent)

	mw.botInfo = &BotInfo{
		Name: connectedEvt.Info.User.Name,
		ID:   connectedEvt.Info.User.ID,
	}
}

func (mw *middleware) handleMessageEvent(payload interface{}) {
	rawMsg := payload.(*slack.MessageEvent)
	var messageType = MessageTypeUnknowMessage

	if strings.HasPrefix(rawMsg.Channel, "D") {
		if rawMsg.User == mw.botInfo.ID {
			messageType = MessageTypeSelfMessage
		} else {
			messageType = MessageTypeDirectMessage
		}
	} else if strings.HasPrefix(rawMsg.Channel, "C") {
		messageType = MessageTypeChannelMessage
	}

	message := Message{
		Type:      messageType,
		Message:   rawMsg.Text,
		Timestamp: rawMsg.Timestamp,
		Channel:   rawMsg.Channel,
		rtm:       mw.slackRTM,
	}

	mw.eventChannel <- &message
}
