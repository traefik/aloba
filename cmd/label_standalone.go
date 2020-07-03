package cmd

import (
	"context"
	"fmt"

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
		return fmt.Errorf("failed to find PR: %w", err)
	}

	for _, issue := range issues {
		err := l.addLabelsToPR(ctx, issue, rc)
		if err != nil {
			return fmt.Errorf("failed to add label to the PR %d: %w", issue.GetNumber(), err)
		}

		if issue.Milestone == nil {
			pr, _, err := l.client.PullRequests.Get(ctx, l.owner, l.repoName, issue.GetNumber())
			if err != nil {
				return fmt.Errorf("failed to get PR %d: %w", issue.GetNumber(), err)
			}

			err = l.addMilestoneToPR(ctx, pr)
			if err != nil {
				return fmt.Errorf("failed to add milestone on PR %d: %w", pr.GetNumber(), err)
			}
		}
	}

	return nil
}
