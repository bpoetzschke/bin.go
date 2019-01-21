package slack_middleware

import (
	"log"
	"os"
	"strings"

	"github.com/bpoetzschke/bin.go/logger"

	"github.com/nlopes/slack"
)

const (
	MessageTypeDirectMessage = "direct_message"
	MessageTypeSelfMessage   = "self_message"
)

type Middleware interface {
	Connect() <-chan *slack.MessageEvent
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

	eventChannel chan *slack.MessageEvent
	botInfo      *BotInfo
	signalCh     chan os.Signal
	logProvider  *log.Logger
}

func (mw *middleware) init() {
	mw.slackApi = slack.New(mw.slackToken)
	mw.eventChannel = make(chan *slack.MessageEvent)
}

func (mw *middleware) Connect() <-chan *slack.MessageEvent {
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
	msg := payload.(*slack.MessageEvent)

	if strings.HasPrefix(msg.Channel, "D") {
		if msg.User == mw.botInfo.ID {
			msg.Type = MessageTypeSelfMessage
		} else {
			msg.Type = MessageTypeDirectMessage
		}
	}

	mw.eventChannel <- msg
}

func (mw *middleware) shutdown() {
	logger.Debug("Attempting graceful shutdown!")

	mw.slackRTM.Disconnect()
	close(mw.eventChannel)
}
