package slack_middleware

import "github.com/nlopes/slack"

func NewSlackRTM(slackToken string) SlackRTM {
	slackClient := slack.New(slackToken)
	rtm := slackClient.NewRTM()

	return &slackRTM{
		rtm: rtm,
	}
}

type slackRTM struct {
	rtm *slack.RTM
}

func (r slackRTM) NewConnection(token string, options ...slack.Option) *slack.Client {
	return slack.New(token, options...)
}

func (r slackRTM) AddReaction(name string, item slack.ItemRef) error {
	return r.rtm.AddReaction(name, item)
}

func (r slackRTM) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	return r.rtm.PostMessage(channelID, options...)
}

func (r slackRTM) ManageConnection() {
	r.rtm.ManageConnection()
}

func (r slackRTM) IncomingEvents() chan slack.RTMEvent {
	return r.rtm.IncomingEvents
}
