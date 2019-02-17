package cmd

import (
	"context"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/label"
	"github.com/containous/aloba/milestone"
	"github.com/containous/aloba/options"
	"github.com/google/go-github/github"
)

// RulesConfiguration Stale issues rules configuration
type RulesConfiguration struct {
	Rules  []label.Rule
	Limits label.Limits
}

// Label add labels to Pull Request
func Label(options *options.Label) error {
	if options.Debug {
		log.Println(options)
	}

	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, options.GitHub.Token)

	rc := &RulesConfiguration{}
	meta, err := toml.DecodeFile(options.RulesFilePath, rc)
	if err != nil {
		return err
	}

	if options.Debug {
		log.Printf("Rules: %+v\n", meta)
	}

	if options.WebHook == nil {
		return runStandalone(ctx, client, options.GitHub.Owner, options.GitHub.RepositoryName, rc, options.DryRun)
	}

	return runWebHook(ctx, client, options.GitHub.Owner, options.GitHub.RepositoryName, rc, options.WebHook, options.DryRun)
}

func addMilestoneToPR(ctx context.Context, client *github.Client, owner, repoName string, pr *github.PullRequest) error {
	meta, err := milestone.Detect(ctx, client, owner, repoName, pr)
	if err != nil {
		return err
	}

	if pr.Milestone == nil && meta != nil {
		ir := &github.IssueRequest{
			Milestone: github.Int(meta.ID),
		}
		_, _, errMil := client.Issues.Edit(ctx, owner, repoName, pr.GetNumber(), ir)
		if errMil != nil {
			return errMil
		}
	}
	return nil
}

func addLabelsToPR(ctx context.Context, client *github.Client, owner string, repositoryName string, issue github.Issue, rc *RulesConfiguration, dryRun bool) error {

	var labels []string

	// AREA
	areas, err := label.DetectAreas(ctx, client, owner, repositoryName, issue.GetNumber(), rc.Rules)
	if err != nil {
		return err
	}
	labels = append(labels, areas...)

	// SIZE
	sizeLabel, err := getSizeLabel(ctx, client, owner, repositoryName, issue, rc.Limits)
	if err != nil {
		return err
	}
	if sizeLabel != "" {
		labels = append(labels, sizeLabel)
	}

	// DIFF
	var addedLabels []string
	for _, lb := range labels {
		if !label.HasLabel(issue.Labels, lb) {
			addedLabels = append(addedLabels, lb)
		}
	}

	// STATUS
	if !label.HasStatus(issue.Labels) {
		addedLabels = append(addedLabels, label.StatusNeedsTriage)
	}

	if len(addedLabels) == 0 {
		log.Printf("#%d: No new labels", issue.GetNumber())
	} else {
		if dryRun {
			log.Printf("#%d: %v - %s\n", issue.GetNumber(), addedLabels, issue.GetTitle())
		} else {
			_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repositoryName, issue.GetNumber(), addedLabels)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getSizeLabel(ctx context.Context, client *github.Client, owner string, repositoryName string, issue github.Issue, limits label.Limits) (string, error) {
	size, err := label.GetSizeLabel(ctx, client, owner, repositoryName, issue.GetNumber(), limits)
	if err != nil {
		return "", err
	}
	currentSize := label.GetCurrentSize(issue.Labels)
	if currentSize != size {
		if currentSize != "" {
			// Remove current size
			_, err := client.Issues.RemoveLabelForIssue(ctx, owner, repositoryName, issue.GetNumber(), currentSize)
			if err != nil {
				return "", err
			}
		}
		return size, nil
	}
	return "", nil
}
