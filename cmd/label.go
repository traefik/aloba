package cmd

import (
	"context"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/containous/myrmica-aloba/internal/gh"
	"github.com/containous/myrmica-aloba/internal/search"
	"github.com/containous/myrmica-aloba/label"
	"github.com/google/go-github/github"
)

type RulesConfiguration struct {
	Rules  []label.Rule
	Limits label.Limits
}

// Label add labels to Pull Request
func Label(options *LabelOptions) error {

	if options.Debug {
		log.Println(options)
	}

	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, options.GitHubToken)

	rc := &RulesConfiguration{}
	meta, err := toml.DecodeFile(options.RulesFilePath, rc)
	if err != nil {
		return err
	}

	if options.Debug {
		log.Printf("Rules: %+v\n", meta)
	}

	if options.WebHook {
		return runWebHook(client, ctx, options.Owner, options.RepositoryName, rc, options.DryRun)
	}
	return runStandalone(client, ctx, options.Owner, options.RepositoryName, rc, options.DryRun)
}

func runStandalone(client *github.Client, ctx context.Context, owner string, repositoryName string, rc *RulesConfiguration, dryRun bool) error {
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
		err := addLabelsToPR(client, ctx, owner, repositoryName, issue, rc, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func addLabelsToPR(client *github.Client, ctx context.Context, owner string, repositoryName string, issue github.Issue, rc *RulesConfiguration, dryRun bool) error {

	labels := []string{}

	// AREA
	areas, err := label.DetectAreas(client, ctx, owner, repositoryName, issue.GetNumber(), rc.Rules)
	if err != nil {
		return err
	}
	labels = append(labels, areas...)

	// SIZE
	sizeLabel, err := getSizeLabel(client, ctx, owner, repositoryName, issue, rc.Limits)
	if err != nil {
		return err
	}
	if sizeLabel != "" {
		labels = append(labels, sizeLabel)
	}

	// DIFF
	addedLabels := []string{}
	for _, lb := range labels {
		if !label.ExistsLabel(issue.Labels, lb) {
			addedLabels = append(addedLabels, lb)
		}
	}

	// STATUS
	if !label.HasStatus(issue.Labels) {
		addedLabels = append(addedLabels, label.StatusNeedsTriage)
	}

	if len(addedLabels) == 0 {
		log.Println("no new labels")
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

func getSizeLabel(client *github.Client, ctx context.Context, owner string, repositoryName string, issue github.Issue, limits label.Limits) (string, error) {
	size, err := label.GetSizeLabel(client, ctx, owner, repositoryName, issue.GetNumber(), limits)
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
