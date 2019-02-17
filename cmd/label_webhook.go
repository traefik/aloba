package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/label"
	"github.com/containous/aloba/options"
	"github.com/google/go-github/github"
	ghw "github.com/ldez/ghwebhook"
	"github.com/ldez/ghwebhook/eventtype"
)

func runWebHook(ctx context.Context, client *github.Client, owner string, repositoryName string, rc *RulesConfiguration, opts *options.WebHook, dryRun bool) error {
	handlers := ghw.NewEventHandlers().
		OnPullRequest(onPullRequest(ctx, client, owner, repositoryName, rc, dryRun)).
		OnPullRequestReview(onPullRequestReview(ctx, client, owner, repositoryName, dryRun)).
		OnIssues(onIssue(ctx, client, owner, repositoryName, dryRun))

	hook := ghw.NewWebHook(handlers,
		ghw.WithPort(opts.Port),
		ghw.WithSecret(opts.Secret),
		ghw.WithEventTypes(eventtype.PullRequest, eventtype.PullRequestReview, eventtype.Issues))
	return hook.ListenAndServe()
}

func onIssue(ctx context.Context, client *github.Client, owner string, repositoryName string, dryRun bool) func(payload *github.WebHookPayload, event *github.IssuesEvent) {
	return func(_ *github.WebHookPayload, event *github.IssuesEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.IssuesEvent) {
				err := onIssueOpened(ctx, client, event, owner, repositoryName, dryRun)
				if err != nil {
					log.Println(err)
				}
			}(event)
		}
	}
}

func onPullRequest(ctx context.Context, client *github.Client, owner string, repositoryName string, rc *RulesConfiguration, dryRun bool) func(*github.WebHookPayload, *github.PullRequestEvent) {
	return func(_ *github.WebHookPayload, event *github.PullRequestEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.PullRequestEvent) {
				err := onPullRequestOpened(ctx, client, event, owner, repositoryName, rc, dryRun)
				if err != nil {
					log.Println(err)
				}
			}(event)
		}
	}
}

func onPullRequestReview(ctx context.Context, client *github.Client, owner string, repositoryName string, dryRun bool) func(*github.WebHookPayload, *github.PullRequestReviewEvent) {
	return func(_ *github.WebHookPayload, event *github.PullRequestReviewEvent) {
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
