package slack_middleware

import (
	"log"
	"strings"

	"github.com/nlopes/slack"

	"github.com/bpoetzschke/bin.go/logger"
)

type Middleware interface {
	Connect() <-chan *IncomingMessage
	GetBotInfo() *BotInfo
	PostMessage(msg OutgoingMessage) error
}

func NewMiddleware(slackToken string) Middleware {
	mw := middleware{slackToken: slackToken}
	mw.init()

	return &mw
}

type middleware struct {
	slackToken string

	slackApi *slack.Client
	slackRTM SlackRTM

	eventChannel chan *IncomingMessage
	botInfo      *BotInfo
	logProvider  *log.Logger
}

func (mw *middleware) init() {
	mw.slackRTM = NewSlackRTM(mw.slackToken)
	mw.eventChannel = make(chan *IncomingMessage, 1)
}

func (mw *middleware) Connect() <-chan *IncomingMessage {
	go mw.slackRTM.ManageConnection()

	go func() {
		for {
			select {
			case evt := <-mw.slackRTM.IncomingEvents():
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

func (mw *middleware) PostMessage(msg OutgoingMessage) error {
	msgOptions := []slack.MsgOption{
		slack.MsgOptionText(msg.Message, false),
		slack.MsgOptionUsername("bingo"),
		slack.MsgOptionAsUser(true),
	}

	attachments := []slack.Attachment{}

	for _, attachment := range msg.Attachments {
		attachments = append(attachments, slack.Attachment{
			ImageURL: attachment,
		})
	}

	msgOptions = append(msgOptions, slack.MsgOptionAttachments(attachments...))

	_, _, err := mw.slackRTM.PostMessage(msg.Channel, msgOptions...)

	return err
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

	if mw.botInfo == nil {
		logger.Warning("Could not determine whether message came from bot user.")
	} else if rawMsg.User == mw.botInfo.ID {
		logger.Debug("Skip handling message since it came from bot user.")
		return
	}

	var messageType = MessageTypeUnknownMessage

	if strings.HasPrefix(rawMsg.Channel, "D") {
		if rawMsg.User == mw.botInfo.ID {
			messageType = MessageTypeSelfMessage
		} else {
			messageType = MessageTypeDirectMessage
		}
	} else if strings.HasPrefix(rawMsg.Channel, "C") {
		messageType = MessageTypeChannelMessage
	}

	message := IncomingMessage{
		BaseMessage: BaseMessage{
			Message: rawMsg.Text,
			Channel: rawMsg.Channel,
		},
		Type:      messageType,
		Timestamp: rawMsg.Timestamp,
		UserID:    rawMsg.User,
		rtm:       mw.slackRTM,
	}

	mw.eventChannel <- &message
}
