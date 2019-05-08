package slack_middleware

import "github.com/nlopes/slack"

func (m Message) React(reactions ...string) error {
	for _, reaction := range reactions {
		if err := m.rtm.AddReaction(reaction, slack.NewRefToMessage(m.Channel, m.Timestamp)); err != nil {
			return err
		}
	}
	return nil
}
