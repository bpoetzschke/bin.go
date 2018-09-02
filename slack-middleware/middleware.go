package slack_middleware

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nlopes/slack"
)

const (
	MessageTypeDirectMessage = "direct_message"
	MessageTypeSelfMessage   = "self_message"
)

type Middleware interface {
	Connect() <-chan *slack.MessageEvent
	GetBotInfo() *BotInfo
	SetLogger(*log.Logger)
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

	// setup handling for graceful shutdown
	mw.eventChannel = make(chan *slack.MessageEvent, 1)
	mw.signalCh = make(chan os.Signal, 1)
	signal.Notify(mw.signalCh, os.Interrupt, syscall.SIGTERM)
}

func (mw *middleware) Connect() <-chan *slack.MessageEvent {
	mw.slackRTM = mw.slackApi.NewRTM()
	go mw.slackRTM.ManageConnection()

	go func() {
		shutdownInitiated := false
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

			case <-mw.signalCh:
				// handle repetitive SIGTERM event
				if shutdownInitiated {
					close(mw.eventChannel)
					return
				}
				shutdownInitiated = true
				mw.shutdown()

			default:
			}
		}
	}()

	return mw.eventChannel
}

func (mw *middleware) GetBotInfo() *BotInfo {
	return mw.botInfo
}

func (mw *middleware) SetLogger(logProvider *log.Logger) {
	if logProvider != nil {
		slack.SetLogger(logProvider)
		mw.slackApi.SetDebug(true)
	} else {
		mw.slackApi.SetDebug(false)
	}
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
	mw.log("Attempting graceful shutdown!")

	mw.slackRTM.Disconnect()
	close(mw.eventChannel)
}

func (mw *middleware) log(msg string) {
	if mw.logProvider != nil {
		mw.logProvider.Print(msg)
	}
}
