package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/label"
	"github.com/google/go-github/github"
	ghw "github.com/ldez/ghwebhook"
	"github.com/ldez/ghwebhook/eventtype"
)

func runWebHook(ctx context.Context, client *github.Client, owner string, repositoryName string, rc *RulesConfiguration, dryRun bool) error {
	handlers := ghw.NewEventHandlers().
		OnPullRequest(onPullRequest(ctx, client, owner, repositoryName, rc, dryRun)).
		OnPullRequestReview(onPullRequestReview(ctx, client, owner, repositoryName, dryRun))

	hook := ghw.NewWebHook(handlers, ghw.WithPort(5000), ghw.WithEventTypes(eventtype.PullRequest, eventtype.PullRequestReview))
	err := hook.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func onPullRequest(ctx context.Context, client *github.Client, owner string, repositoryName string, rc *RulesConfiguration, dryRun bool) func(*github.WebHookPayload, *github.PullRequestEvent) {
	return func(payload *github.WebHookPayload, event *github.PullRequestEvent) {
		if event.GetAction() == "opened" {
			go func(event *github.PullRequestEvent) {
				issue, _, err := client.Issues.Get(ctx, owner, repositoryName, event.GetNumber())
				if err != nil {
					log.Println(err)
					return
				}

				err = addLabelsToPR(ctx, client, owner, repositoryName, *issue, rc, dryRun)
				if err != nil {
					log.Println(err)
				}
			}(event)
		}
	}
}

func onPullRequestReview(ctx context.Context, client *github.Client, owner string, repositoryName string, dryRun bool) func(*github.WebHookPayload, *github.PullRequestReviewEvent) {
	return func(payload *github.WebHookPayload, event *github.PullRequestReviewEvent) {
		if event.GetAction() == "submitted" {
			if strings.ToUpper(event.Review.GetState()) == gh.ChangesRequested {
				go func(event *github.PullRequestReviewEvent) {

					issue, _, err := client.Issues.Get(ctx, owner, repositoryName, event.PullRequest.GetNumber())
					if err != nil {
						log.Println(err)
						return
					}

					if label.HasLabel(issue.Labels, label.ContributorWaitingForCorrections) {
						return
					}

					if dryRun {
						log.Printf("#%d: Add %v\n", issue.GetNumber(), label.ContributorWaitingForCorrections)
					} else {
						_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repositoryName, issue.GetNumber(), []string{label.ContributorWaitingForCorrections})
						if err != nil {
							log.Println(err)
							return
						}
					}
				}(event)
			} else if strings.ToUpper(event.Review.GetState()) == gh.Approved {
				go func(event *github.PullRequestReviewEvent) {
					issue, _, err := client.Issues.Get(ctx, owner, repositoryName, event.PullRequest.GetNumber())
					if err != nil {
						log.Println(err)
						return
					}

					err = removeLabel(ctx, client, owner, repositoryName, issue, label.ContributorWaitingForCorrections, dryRun)
					if err != nil {
						log.Println(err)
						return
					}
				}(event)
			}
		}
	}
}

func removeLabel(ctx context.Context, client *github.Client, owner string, repositoryName string, issue *github.Issue, labelName string, dryRun bool) error {
	if label.HasLabel(issue.Labels, labelName) {
		if dryRun {
			log.Printf("#%d: Remove %v\n", issue.GetNumber(), labelName)
		} else {
			_, err := client.Issues.RemoveLabelForIssue(ctx, owner, repositoryName, issue.GetNumber(), labelName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
