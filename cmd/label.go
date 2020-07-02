package cmd

import (
	"context"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/label"
	"github.com/containous/aloba/milestone"
	"github.com/containous/aloba/options"
	"github.com/google/go-github/v27/github"
)

const stateOpened = "opened"

// RulesConfiguration Stale issues rules configuration.
type RulesConfiguration struct {
	Rules  []label.Rule
	Limits label.Limits
}

// Label adds labels to Pull Request.
func Label(options *options.Label) error {
	if options.Debug {
		log.Println(options)
	}

	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, options.GitHub.Token)

	labeler := NewLabeler(client)

	rc := &RulesConfiguration{}
	meta, err := toml.DecodeFile(options.RulesFilePath, rc)
	if err != nil {
		return err
	}

	if options.Debug {
		log.Printf("Rules: %+v\n", meta)
	}

	if options.WebHook == nil {
		return labeler.runStandalone(ctx, options.GitHub.Owner, options.GitHub.RepositoryName, rc, options.DryRun)
	}

	return labeler.runWebHook(ctx, options.GitHub.Owner, options.GitHub.RepositoryName, rc, options.WebHook, options.DryRun)
}

// Labeler adds labels to Pull Request.
type Labeler struct {
	client *github.Client
}

// NewLabeler creates a new Labeler.
func NewLabeler(client *github.Client) *Labeler {
	return &Labeler{client: client}
}

func (l *Labeler) addMilestoneToPR(ctx context.Context, owner, repoName string, pr *github.PullRequest) error {
	meta, err := milestone.Detect(ctx, l.client, owner, repoName, pr)
	if err != nil {
		return err
	}

	if pr.Milestone == nil && meta != nil {
		ir := &github.IssueRequest{
			Milestone: github.Int(meta.ID),
		}
		_, _, errMil := l.client.Issues.Edit(ctx, owner, repoName, pr.GetNumber(), ir)
		if errMil != nil {
			return errMil
		}
	}
	return nil
}

func (l *Labeler) addLabelsToPR(ctx context.Context, owner string, repositoryName string, issue github.Issue, rc *RulesConfiguration, dryRun bool) error {
	var labels []string

	// AREA
	areas, err := label.DetectAreas(ctx, l.client, owner, repositoryName, issue.GetNumber(), rc.Rules)
	if err != nil {
		return err
	}
	labels = append(labels, areas...)

	// SIZE
	sizeLabel, err := l.getSizeLabel(ctx, owner, repositoryName, issue, rc.Limits)
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
		return nil
	}

	if dryRun {
		log.Printf("#%d: %v - %s\n", issue.GetNumber(), addedLabels, issue.GetTitle())
		return nil
	}

	_, _, err = l.client.Issues.AddLabelsToIssue(ctx, owner, repositoryName, issue.GetNumber(), addedLabels)
	if err != nil {
		return err
	}

	return nil
}

func (l *Labeler) getSizeLabel(ctx context.Context, owner string, repositoryName string, issue github.Issue, limits label.Limits) (string, error) {
	size, err := label.GetSizeLabel(ctx, l.client, owner, repositoryName, issue.GetNumber(), limits)
	if err != nil {
		return "", err
	}
	currentSize := label.GetCurrentSize(issue.Labels)
	if currentSize != size {
		if currentSize != "" {
			// Remove current size
			_, err := l.client.Issues.RemoveLabelForIssue(ctx, owner, repositoryName, issue.GetNumber(), currentSize)
			if err != nil {
				return "", err
			}
		}
		return size, nil
	}
	return "", nil
}

func (l *Labeler) onIssueOpened(ctx context.Context, event *github.IssuesEvent, owner, repositoryName string, dryRun bool) error {
	// add sleep due to some GitHub latency
	time.Sleep(1 * time.Second)

	issue, _, err := l.client.Issues.Get(ctx, owner, repositoryName, event.Issue.GetNumber())
	if err != nil {
		return err
	}

	if len(issue.Labels) == 0 {
		if dryRun {
			log.Printf("Add %q label to %d", label.StatusNeedsTriage, event.Issue.GetNumber())
			return nil
		}

		_, _, err = l.client.Issues.AddLabelsToIssue(ctx, owner, repositoryName, issue.GetNumber(), []string{label.StatusNeedsTriage})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Labeler) onPullRequestOpened(ctx context.Context, event *github.PullRequestEvent, owner, repositoryName string, rc *RulesConfiguration, dryRun bool) error {
	// add sleep due to some GitHub latency
	time.Sleep(1 * time.Second)

	issue, _, err := l.client.Issues.Get(ctx, owner, repositoryName, event.GetNumber())
	if err != nil {
		return err
	}

	err = l.addLabelsToPR(ctx, owner, repositoryName, *issue, rc, dryRun)
	if err != nil {
		return err
	}

	if event.PullRequest.Milestone == nil {
		err = l.addMilestoneToPR(ctx, owner, repositoryName, event.PullRequest)
		if err != nil {
			return err
		}
	}

	return nil
}
