package cmd

import (
	"context"

	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/report"
)

// Report create a report and publish on Slack
func Report(options *ReportOptions) error {
	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, options.GitHubToken)

	model, err := report.MakeReport(client, ctx, options.Owner, options.RepositoryName)
	if err != nil {
		return err
	}

	if options.Debug || options.DryRun {
		report.DisplayReport(model)
	}

	if options.DryRun {
		return nil
	}
	return report.SendToSlack(options.SlackToken, options.ChannelID, options.IconEmoji, options.BotName, model)
}
