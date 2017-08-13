package report

import (
	"strings"

	"github.com/containous/aloba/options"
	"github.com/nlopes/slack"
)

// SendToSlack Create and send a report to Slack
func SendToSlack(options *options.Slack, model *Model) error {

	msgParts := []string{}
	if len(model.withReviews) != 0 {
		msgParts = append(msgParts, "With reviews:", makeMessage(model.withReviews, false))
	}
	if len(model.noReviews) != 0 {
		msgParts = append(msgParts, "No reviews:", makeMessage(model.noReviews, false))
	}
	if len(model.designReview) != 0 {
		msgParts = append(msgParts, "Need design review:", makeMessage(model.designReview, false))
	}

	if len(msgParts) != 0 {
		message := strings.Join(msgParts, "\n")

		api := slack.New(options.Token)

		ppm := slack.PostMessageParameters{
			IconEmoji: options.IconEmoji,
			Username:  options.BotName,
		}

		_, _, err := api.PostMessage(options.ChannelID, message, ppm)
		if err != nil {
			return err
		}
	}
	return nil
}
