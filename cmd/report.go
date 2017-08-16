package cmd

import (
	"context"

	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/options"
	"github.com/containous/aloba/report"
)

// Report create a report and publish on Slack
func Report(options *options.Report) error {
	if options.ServerMode {
		server := &server{options: options}
		return server.ListenAndServe()
	}

	return launch(options)
}

func launch(options *options.Report) error {
	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, options.GitHub.Token)

	model, err := report.MakeReport(ctx, client, options.GitHub.Owner, options.GitHub.RepositoryName)
	if err != nil {
		return err
	}

	if options.Debug || options.DryRun {
		report.DisplayReport(model)
	}

	if options.DryRun {
		return nil
	}
	return report.SendToSlack(options.Slack, model)
}
