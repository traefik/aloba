package report

import (
	"strings"

	"github.com/containous/aloba/options"
	"github.com/nlopes/slack"
)

// SendToSlack Create and send a report to Slack
func SendToSlack(options *options.Slack, model *Model) error {

	message := makeSlackMessage(model)

	if len(message) != 0 {
		return postSlackMessage(options, message)
	}
	return nil
}

func postSlackMessage(options *options.Slack, message string) error {
	api := slack.New(options.Token)

	ppm := slack.PostMessageParameters{
		IconEmoji: options.IconEmoji,
		Username:  options.BotName,
	}

	_, _, err := api.PostMessage(options.ChannelID,
		slack.MsgOptionPostMessageParameters(ppm),
		slack.MsgOptionText(message, false),
	)
	return err
}

func makeSlackMessage(model *Model) string {
	var msgParts []string

	if len(model.withReviews) != 0 {
		msgParts = append(msgParts, "*With reviews:*", makeMessage(model.withReviews, false))
	}
	if len(model.noReviews) != 0 {
		msgParts = append(msgParts, "*Without review:*", makeMessage(model.noReviews, false))
	}
	if len(model.designReview) != 0 {
		msgParts = append(msgParts, "*Need design review:*", makeMessage(model.designReview, false))
	}

	msg := strings.Join(msgParts, "\n")

	if len(msg) != 0 {
		msg = "By your powers combined I am Captain PR! We gonna take PR down to zero!\n\n" + msg
		//TODO must authorize bot to call @channel
		// msg = "<!channel> By your powers combined I am Captain PR! We gonna take PR down to zero!\n\n" + msg
	}

	return msg
}
