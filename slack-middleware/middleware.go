package slack_middleware

import (
	"fmt"
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
}

/*

func main() {
	cmd := Cmd{Closed: make(chan struct{})}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	log.Println("Listening for signals")

	// Block until one of the signals above is received
	<-signalCh
	log.Println("Signal received, initializing clean shutdown...")
	go cmd.Close()

	// Block again until another signal is received, a shutdown timeout elapses,
	// or the Command is gracefully closed
	log.Println("Waiting for clean shutdown...")
	select {
	case <-signalCh:
		log.Println("second signal received, initializing hard shutdown")
	case <-time.After(time.Second * 30):
		log.Println("time limit reached, initializing hard shutdown")
	case <-cmd.Closed:
		log.Println("server shutdown completed")
	}

	// goodbye.
}

*/

func (mw *middleware) init() {
	mw.slackApi = slack.New(mw.slackToken)

	//logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	//slack.SetLogger(logger)
	//mw.slackApi.SetDebug(true)

	mw.eventChannel = make(chan *slack.MessageEvent, 1)
	mw.signalCh = make(chan os.Signal, 1)
	signal.Notify(mw.signalCh, os.Interrupt, syscall.SIGTERM)
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

			case <-mw.signalCh:
				close(mw.eventChannel)
				mw.shutdown()
				return

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
	fmt.Println("Graceful shutdown!")

	mw.slackRTM.Disconnect()
}
