package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/v27/github"
	"github.com/traefik/aloba/internal/gh"
	"github.com/traefik/aloba/label"
	"github.com/traefik/aloba/milestone"
	"github.com/traefik/aloba/options"
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

	labeler := NewLabeler(client, options.GitHub.Owner, options.GitHub.RepositoryName)
	labeler.DryRun = options.DryRun

	rc := &RulesConfiguration{}
	meta, err := toml.DecodeFile(options.RulesFilePath, rc)
	if err != nil {
		return fmt.Errorf("failed to read TOML file (%s): %w", options.RulesFilePath, err)
	}

	if options.Debug {
		log.Printf("Rules: %+v\n", meta)
	}

	if options.WebHook == nil {
		return labeler.runStandalone(ctx, rc)
	}

	return labeler.runWebHook(ctx, rc, options.WebHook)
}

// Labeler adds labels to Pull Request.
type Labeler struct {
	client   *github.Client
	owner    string
	repoName string
	DryRun   bool
}

// NewLabeler creates a new Labeler.
func NewLabeler(client *github.Client, owner, repoName string) *Labeler {
	return &Labeler{
		client:   client,
		owner:    owner,
		repoName: repoName,
	}
}

func (l *Labeler) addMilestoneToPR(ctx context.Context, pr *github.PullRequest) error {
	meta, err := milestone.Detect(ctx, l.client, l.owner, l.repoName, pr)
	if err != nil {
		return fmt.Errorf("failed to detect milestone for PR %d: %w", pr.GetNumber(), err)
	}

	if pr.Milestone == nil && meta != nil {
		ir := &github.IssueRequest{
			Milestone: github.Int(meta.ID),
		}
		_, _, errMil := l.client.Issues.Edit(ctx, l.owner, l.repoName, pr.GetNumber(), ir)
		if errMil != nil {
			return fmt.Errorf("failed to add milestone to PR %d: %w", pr.GetNumber(), errMil)
		}
	}
	return nil
}

func (l *Labeler) addLabelsToPR(ctx context.Context, issue github.Issue, rc *RulesConfiguration) error {
	var labels []string

	// AREA
	areas, err := label.DetectAreas(ctx, l.client, l.owner, l.repoName, issue.GetNumber(), rc.Rules)
	if err != nil {
		return fmt.Errorf("failed to detect area for PR %d: %w", issue.GetNumber(), err)
	}
	labels = append(labels, areas...)

	// SIZE
	sizeLabel, err := l.getSizeLabel(ctx, issue, rc.Limits)
	if err != nil {
		return fmt.Errorf("failed to computethe size of the PR %d: %w", issue.GetNumber(), err)
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

	if l.DryRun {
		log.Printf("#%d: %v - %s\n", issue.GetNumber(), addedLabels, issue.GetTitle())
		return nil
	}

	_, _, err = l.client.Issues.AddLabelsToIssue(ctx, l.owner, l.repoName, issue.GetNumber(), addedLabels)
	if err != nil {
		return fmt.Errorf("failed to add label to PR %d: %w", issue.GetNumber(), err)
	}

	return nil
}

func (l *Labeler) getSizeLabel(ctx context.Context, issue github.Issue, limits label.Limits) (string, error) {
	size, err := label.GetSizeLabel(ctx, l.client, l.owner, l.repoName, issue.GetNumber(), limits)
	if err != nil {
		return "", err
	}

	currentSize := label.GetCurrentSize(issue.Labels)
	if currentSize != size {
		if currentSize != "" {
			// Remove current size
			_, err := l.client.Issues.RemoveLabelForIssue(ctx, l.owner, l.repoName, issue.GetNumber(), currentSize)
			if err != nil {
				return "", fmt.Errorf("failed to remove size label to PR %d: %w", issue.GetNumber(), err)
			}
		}
		return size, nil
	}

	return "", nil
}

func (l *Labeler) onIssueOpened(ctx context.Context, event *github.IssuesEvent) error {
	// add sleep due to some GitHub latency
	time.Sleep(1 * time.Second)

	issue, _, err := l.client.Issues.Get(ctx, l.owner, l.repoName, event.Issue.GetNumber())
	if err != nil {
		return fmt.Errorf("failed to get issue %d: %w", event.Issue.GetNumber(), err)
	}

	if len(issue.Labels) == 0 {
		if l.DryRun {
			log.Printf("Add %q label to %d", label.StatusNeedsTriage, event.Issue.GetNumber())
			return nil
		}

		_, _, err = l.client.Issues.AddLabelsToIssue(ctx, l.owner, l.repoName, issue.GetNumber(), []string{label.StatusNeedsTriage})
		if err != nil {
			return fmt.Errorf("failed to add label to issue %d: %w", issue.GetNumber(), err)
		}
	}

	return nil
}

func (l *Labeler) onPullRequestOpened(ctx context.Context, event *github.PullRequestEvent, rc *RulesConfiguration) error {
	// add sleep due to some GitHub latency
	time.Sleep(1 * time.Second)

	issue, _, err := l.client.Issues.Get(ctx, l.owner, l.repoName, event.GetNumber())
	if err != nil {
		return fmt.Errorf("failed to get PR %d: %w", event.GetNumber(), err)
	}

	err = l.addLabelsToPR(ctx, *issue, rc)
	if err != nil {
		return fmt.Errorf("failed to add labels on PR %d: %w", issue.GetNumber(), err)
	}

	if event.PullRequest.Milestone == nil {
		err = l.addMilestoneToPR(ctx, event.PullRequest)
		if err != nil {
			return fmt.Errorf("failed to add milestone on PR %d: %w", issue.GetNumber(), err)
		}
	}

	return nil
}
