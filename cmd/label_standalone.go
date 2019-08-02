package cmd

import (
	"context"

	"github.com/containous/aloba/internal/search"
	"github.com/containous/aloba/label"
	"github.com/google/go-github/v27/github"
)

func runStandalone(ctx context.Context, client *github.Client, owner string, repositoryName string, rc *RulesConfiguration, dryRun bool) error {
	issues, err := search.FindOpenPR(ctx, client, owner, repositoryName,
		search.WithExcludedLabels(
			label.SizeLabelPrefix+label.Small,
			label.SizeLabelPrefix+label.Medium,
			label.SizeLabelPrefix+label.Large,
			label.WIP))
	if err != nil {
		return err
	}

	for _, issue := range issues {
		err := addLabelsToPR(ctx, client, owner, repositoryName, issue, rc, dryRun)
		if err != nil {
			return err
		}
		if issue.Milestone == nil {
			pr, _, err := client.PullRequests.Get(ctx, owner, repositoryName, issue.GetNumber())
			if err != nil {
				return err
			}
			err = addMilestoneToPR(ctx, client, owner, repositoryName, pr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
