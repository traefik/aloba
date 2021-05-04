package cmd

import (
	"context"

	"github.com/traefik/aloba/internal/gh"
	"github.com/traefik/aloba/options"
	"github.com/traefik/aloba/report"
)

// Report create a report and publish on Slack.
func Report(options *options.Report) error {
	if options.ServerMode {
		server := &server{options: options}
		return server.ListenAndServe()
	}

	return launch(options)
}

func launch(options *options.Report) error {
	ctx := context.Background()

	reporter := report.NewReporter(gh.NewGitHubClient(ctx, options.GitHub.Token), options.GitHub.Owner, options.GitHub.RepositoryName)

	model, err := reporter.Make(ctx)
	if err != nil {
		return err
	}

	if options.DryRun {
		report.DisplayReport(model)
		return nil
	}

	return report.SendToSlack(options.Slack, model)
}
