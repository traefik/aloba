package cmd

import (
	"context"

	"github.com/containous/aloba/internal/search"
	"github.com/containous/aloba/label"
)

func (l *Labeler) runStandalone(ctx context.Context, owner string, repositoryName string, rc *RulesConfiguration, dryRun bool) error {
	issues, err := search.FindOpenPR(ctx, l.client, owner, repositoryName,
		search.WithExcludedLabels(
			label.SizeLabelPrefix+label.Small,
			label.SizeLabelPrefix+label.Medium,
			label.SizeLabelPrefix+label.Large,
			label.WIP))
	if err != nil {
		return err
	}

	for _, issue := range issues {
		err := l.addLabelsToPR(ctx, owner, repositoryName, issue, rc, dryRun)
		if err != nil {
			return err
		}

		if issue.Milestone == nil {
			pr, _, err := l.client.PullRequests.Get(ctx, owner, repositoryName, issue.GetNumber())
			if err != nil {
				return err
			}

			err = l.addMilestoneToPR(ctx, owner, repositoryName, pr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
