package report

import (
	"strings"

	"github.com/nlopes/slack"
)

func SendToSlack(slackToken string, channelID string, iconEmoji string, botName string, model *model) error {

	api := slack.New(slackToken)

	ppm := slack.PostMessageParameters{
		IconEmoji: iconEmoji,
		Username:  botName,
	}

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

		_, _, err := api.PostMessage(channelID, message, ppm)
		if err != nil {
			return err
		}
	}
	return nil
}
