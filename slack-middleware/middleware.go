package slack_middleware

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

const (
	MessageTypeDirectMessage = "direct_message"
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
	slackToken   string
	slackApi     *slack.Client
	eventChannel chan *slack.MessageEvent
	botInfo      *BotInfo
}

func (mw *middleware) init() {
	mw.slackApi = slack.New(mw.slackToken)

	//logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	//slack.SetLogger(logger)
	//mw.slackApi.SetDebug(true)

	mw.eventChannel = make(chan *slack.MessageEvent, 1)
}

func (mw *middleware) Connect() <-chan *slack.MessageEvent {
	rtm := mw.slackApi.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case evt := <-rtm.IncomingEvents:
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
		msg.Type = MessageTypeDirectMessage
	}

	fmt.Printf("Message received: %+v", msg)
}
