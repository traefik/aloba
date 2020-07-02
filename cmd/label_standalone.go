package cmd

import (
	"context"

	"github.com/containous/aloba/internal/search"
	"github.com/containous/aloba/label"
)

func (l *Labeler) runStandalone(ctx context.Context, rc *RulesConfiguration) error {
	issues, err := search.FindOpenPR(ctx, l.client, l.owner, l.repoName,
		search.WithExcludedLabels(
			label.SizeLabelPrefix+label.Small,
			label.SizeLabelPrefix+label.Medium,
			label.SizeLabelPrefix+label.Large,
			label.WIP))
	if err != nil {
		return err
	}

	for _, issue := range issues {
		err := l.addLabelsToPR(ctx, issue, rc)
		if err != nil {
			return err
		}

		if issue.Milestone == nil {
			pr, _, err := l.client.PullRequests.Get(ctx, l.owner, l.repoName, issue.GetNumber())
			if err != nil {
				return err
			}

			err = l.addMilestoneToPR(ctx, pr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
